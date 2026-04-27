package api

import (
	"database/sql"
	"errors"
	"io"
	"logger"
	"mime/multipart"
	"models"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"server/config"
	"server/database"
	"server/internal/exploit"
	"server/websockets"

	json "github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v3"
)

const flagCheckerHostNotConfigureWarnMessage = "Flagchecker host not configured"

// ---------- GET ----------------

// HandleGetConfig returns the current shared configuration of the server.
//
// @Summary Get shared config
// @Description Returns current shared configuration used by clients/exploit runtime.// @Tags config
// @Produce json
// @Security CookieAuth
// @Success 200 {object} ResponseSharedConfig
// @Router /config [get]
func (h *Handler) HandleGetConfig(c fiber.Ctx) error {
	return c.JSON(h.config.GetShared())
}

// HandleGetConfig returns the current full configuration of the server.
//
// @Summary Get full config
// @Description Returns current server and shared configuration.
// @Tags config
// @Produce json
// @Security CookieAuth
// @Success 200 {object} ResponseSharedConfig
// @Router /config [get]
func (h *Handler) HandleGetFullConfig(c fiber.Ctx) error {
	return c.JSON(h.config.GetFullConfig())
}

// HandleGetAllFlags retrieves and returns all the stored flags.
//
// @Summary List flags
// @Description Returns all stored flags.
// @Tags flags
// @Produce json
// @Security CookieAuth
// @Success 200 {object} ResponseFlags
// @Failure 500 {object} ResponseError
// @Router /flags [get]
func (h *Handler) HandleGetAllFlags(c fiber.Ctx) error {
	flags, err := h.store.Queries.GetAllFlags(c.RequestCtx())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch all flags")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{Error: err.Error()})
	}
	if flags == nil {
		flags = []database.Flag{}
	}
	data := ResponseFlags{
		Nflags: int64(len(flags)),
		Flags:  flags,
	}
	return c.JSON(data)
}

// HandleGetStats returns statistics about the server state.
// Returns aggregated statistics for flags. The response body contains a JSON object
// with the key "flags_stats" that maps to an array of aggregated rows (grouped
// statistics returned by the database query). Each element typically contains
// fields such as service name, status and count depending on the underlying
// query implementation.
//
// @Summary Get stats
// @Description Returns aggregated server/flags statistics (grouped counts by service/status).
// @Tags stats
// @Produce json
// @Security CookieAuth
// @Success 200 {object} map[string]any "{"flags_stats": [{"service":"svc","status":1,"count":10}, ...] }"
// @Failure 500 {object} ResponseError
// @Router /stats [get]
func (h *Handler) HandleGetStats(c fiber.Ctx) error {
	rows, err := h.store.Queries.FlagsStats(c.RequestCtx())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch flags stats")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{Error: err.Error()})
	}

	if rows == nil {
		return c.JSON(fiber.Map{"flags_stats": []any{}})
	}

	return c.JSON(fiber.Map{"flags_stats": rows})
}

// HandleGetPaginatedFlags returns paginated and filtered flags.
//
// @Summary List flags (paginated)
// @Description Returns paginated flags with optional status/team/search filters.
// @Tags flags
// @Produce json
// @Security CookieAuth
// @Param limit path int true "Page size"
// @Param offset query int false "Offset"
// @Param status query int false "Status filter"
// @Param service query string false "Service filter"
// @Param team query string false "Team filter"
// @Param search query string false "Search text"
// @Param search_field query string false "Search field"
// @Success 200 {object} ResponseFlags
// @Failure 500 {object} ResponseError
// @Router /flags/{limit} [get]
func (h *Handler) HandleGetPaginatedFlags(c fiber.Ctx) error {
	limit, err := fiber.Params[int](c, "limit", config.DefaultLimit), error(nil)
	if err != nil || limit <= 0 {
		logger.Log.Warn().Msg("Invalid or missing limit parameter")
		limit = config.DefaultLimit
	}
	offset := fiber.Query[int](c, "offset", config.DefaultOffset)
	if offset < 0 {
		logger.Log.Warn().Msg("Invalid offset parameter, using default")
		offset = config.DefaultOffset
	}

	// Build filter options from query parameters
	optsStatus := int64(fiber.Query[int](c, "status", 5))
	optsService := strings.TrimSpace(c.Query("service", ""))
	teamStr := strings.TrimSpace(c.Query("team", ""))

	var serviceNull sql.NullString
	if optsService != "" {
		serviceNull = sql.NullString{String: optsService, Valid: true}
	} else {
		serviceNull = sql.NullString{Valid: false}
	}

	// Parse team ID if provided
	var teamID uint16
	if teamStr != "" {
		if parsed, err := strconv.ParseUint(teamStr, 10, 16); err == nil {
			teamID = uint16(parsed)
		} else {
			logger.Log.Warn().Err(err).Msg("Invalid team parameter, ignoring")
		}
	}

	opts := database.GetFilteredFlagsParams{
		Status: sql.NullInt64{Int64: optsStatus, Valid: optsStatus != 5}, // Simple filter for the status (UNSUBMITTED/ACCEPTED/DENIED/ERROR)
		TeamID: sql.NullInt64{Int64: int64(teamID), Valid: teamID != 0},  // Filter by team ID (0 means not provided)
		Search: sql.NullString{
			String: strings.TrimSpace(c.Query("search", "")),
			Valid:  strings.TrimSpace(c.Query("search", "")) != "",
		}, // Value of the search query
		SearchField: sql.NullString{
			String: strings.TrimSpace(c.Query("search_field", "flag_code")),
			Valid:  true,
		}, // Field to apply the search query to (default: flag_code)
		Limit:  sql.NullInt64{Int64: int64(limit), Valid: true},
		Offset: sql.NullInt64{Int64: int64(offset), Valid: true},
	}

	flags, err := h.store.Queries.GetFilteredFlags(c.RequestCtx(), opts)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch filtered flags")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{Error: err.Error()})
	}

	optsCount := database.CountFilteredFlagsParams{
		Status:      sql.NullInt64{Int64: optsStatus, Valid: optsStatus != 5}, // Simple filter for the status (SUBMITTED/UNSUBMITTED/ACCEPTED/DENIED/ERROR)
		ServiceName: serviceNull,                                              // Filter by service name
		TeamID:      sql.NullInt64{Int64: int64(teamID), Valid: teamID != 0},  // Filter by team ID (0 means not provided)
		Search: sql.NullString{
			String: strings.TrimSpace(c.Query("search", "")),
			Valid:  strings.TrimSpace(c.Query("search", "")) != "",
		}, // Value of the search query
		SearchField: sql.NullString{
			String: strings.TrimSpace(c.Query("search_field", "flag_code")),
			Valid:  true,
		}, // Field to apply the search query (default: flag_code)

	}

	// Get filtered count for accurate pagination
	nFlags, err := h.store.Queries.CountFilteredFlags(c.RequestCtx(), optsCount)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to count filtered flags")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{Error: err.Error()})
	}

	if flags == nil {
		flags = []database.Flag{}
	}

	return c.JSON(ResponseFlags{
		Nflags: nFlags,
		Flags:  flags,
	})
}

// HandleGetFlag retrieves a single flag by its ID.
//
// @Summary List protocols
// @Description Returns protocol plugins available on server.
// @Tags protocols
// @Produce json
// @Success 200 {object} map[string][]string
// @Failure 500 {object} ResponseError
// @Router /protocols [get]
func (*Handler) HandleGetProtocols(c fiber.Ctx) error {
	searchPaths := []string{
		"pkg/protocols",
		"protocols",
	}

	var protocolNames []string
	for _, exploitPath := range searchPaths {
		if protocols, err := os.ReadDir(exploitPath); err == nil {
			for _, entry := range protocols {
				if entry.IsDir() {
					protocolNames = append(protocolNames, strings.Split(entry.Name(), ".")[0])
				} else if matched, _ := filepath.Match("*.so", entry.Name()); matched {
					protocolNames = append(protocolNames, strings.Split(entry.Name(), ".")[0])
				}
			}
			break
		}
	}

	if len(protocolNames) == 0 {
		logger.Log.Error().Msg("Failed to read protocols directory")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{Error: "No protocols found"})
	}

	return c.JSON(fiber.Map{
		"protocols": protocolNames,
	})
}

// ---------- POST ----------------

// HandlePostFlags processes a batch of flags submitted in the request.
//
// @Summary Submit flag batch (store only)
// @Description Stores a batch of flags without immediate external checker submission.
// @Tags submit
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body models.SubmitFlagsRequest true "Flags payload"
// @Success 200 {object} ResponseSuccess
// @Failure 422 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /submit-flags [post]
func (h *Handler) HandlePostFlags(c fiber.Ctx) error {
	var payload models.SubmitFlagsRequest
	if err := c.Bind().Body(&payload); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid SubmitFlags payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ResponseError{Error: err.Error()})
	}
	for _, flag := range payload.Flags {
		if err := h.store.Queries.AddFlag(c.RequestCtx(), database.MapFromFlagToDBParams(flag)); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to insert flags")
			return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{Error: err.Error()})
		}
	}

	return c.JSON(ResponseSuccess{Message: "Flags submitted successfully"})
}

// HandlePostFlag processes a single flag and optionally submits it to an external checker.
//
// @Summary Submit one flag
// @Description Stores one flag and attempts checker submission immediately.
// @Tags submit
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body models.SubmitFlagRequest true "Flag payload"
// @Success 200 {object} ResponseSuccess
// @Failure 422 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Failure 503 {object} ResponseError
// @Router /submit-flag [post]
func (h *Handler) HandlePostFlag(c fiber.Ctx) error {
	var payload models.SubmitFlagRequest
	if err := c.Bind().Body(&payload); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid SubmitFlag payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ResponseError{Error: err.Error()})
	}

	if err := h.store.Queries.AddFlag(c.RequestCtx(), database.MapFromFlagToDBParams(payload.Flag)); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to insert single flag")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{Error: err.Error()})
	}

	if h.config.GetURLFlagChecker() == "" {
		logger.Log.Warn().Msg(flagCheckerHostNotConfigureWarnMessage)
		return c.Status(fiber.StatusServiceUnavailable).JSON(ResponseError{
			Error: flagCheckerHostNotConfigureWarnMessage,
		})
	}

	flags := []string{payload.Flag.FlagCode}
	response, err := config.Submit(h.config.GetURLFlagChecker(), h.config.GetTeamToken(), flags)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Flagchecker submission failed")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Error:   "Failed to submit flag",
			Details: err.Error(),
		})
	}

	h.runner.UpdateFlags(response)

	return c.JSON(ResponseSuccess{Message: "Flag submitted successfully"})
}

// HandlePostFlag processes a single flag and optionally submits it to an external checker.
//
// @Summary Submit flag batch (standalone)
// @Description Stores a batch and attempts checker submission immediately.
// @Tags submit
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body models.SubmitFlagsRequest true "Flags payload"
// @Success 200 {object} ResponseSuccess
// @Failure 422 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Failure 503 {object} ResponseError
// @Router /submit-flags-standalone [post]
func (h *Handler) HandlePostFlagsStandalone(c fiber.Ctx) error {
	var payload models.SubmitFlagsRequest
	if err := c.Bind().Body(&payload); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid SubmitFlag payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ResponseError{Error: err.Error()})
	}

	for _, flag := range payload.Flags {
		if err := h.store.Queries.AddFlag(c.RequestCtx(), database.MapFromFlagToDBParams(flag)); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to insert single flag")
			return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{Error: err.Error()})
		}
	}

	if h.config.GetURLFlagChecker() == "" {
		logger.Log.Warn().Msg(flagCheckerHostNotConfigureWarnMessage)
		return c.Status(fiber.StatusServiceUnavailable).JSON(ResponseError{
			Error: flagCheckerHostNotConfigureWarnMessage,
		})
	}
	flags := make([]string, len(payload.Flags))

	for i, flag := range payload.Flags {
		flags[i] = flag.FlagCode
		if flag.FlagCode == "" {
			logger.Log.Warn().Msg("Empty flag code found, skipping submission")
			continue
		}
	}

	response, err := config.Submit(h.config.GetURLFlagChecker(), h.config.GetTeamToken(), flags)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Flagchecker submission failed")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Error:   "Failed to submit flag",
			Details: err.Error(),
		})
	}

	h.runner.UpdateFlags(response)

	return c.JSON(ResponseSuccess{Message: "Flag submitted successfully"})
}

// HandlePostConfig updates the full configuration and restarts the flag processing loop.
//
// @Summary Update full config
// @Description Updates server/shared configuration and restarts background runner loops.
// @Tags config
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body UpdateConfigRequest true "Full config payload"
// @Success 200 {object} ResponseSuccess
// @Failure 422 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /config [post]
func (h *Handler) HandlePostConfig(c fiber.Ctx) error {
	var payload UpdateConfigRequest
	if err := c.Bind().Body(&payload); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to parse config payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ResponseError{
			Error:   "Invalid config payload",
			Details: err.Error(),
		})
	}

	nextConfig := payload.Config
	nextConfig.Configured = true
	nextConfig.Shared.Configured = true

	h.config.SetFullConfig(nextConfig)
	h.config.SetConfigured(true)

	h.runner.Run()

	cfgJSON, err := json.Marshal(h.config.GetShared())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to marshal config")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Error:   "Failed to marshal config",
			Details: err.Error(),
		})
	}

	event := websockets.Event{
		Type:    websockets.ConfigMessage,
		Payload: cfgJSON,
	}

	if websockets.GlobalManager != nil {
		for client := range websockets.GlobalManager.Clients {
			client.Egress <- event
		}
	}

	return c.JSON(ResponseSuccess{Message: "Configuration updated successfully"})
}

// HandleDeleteFlag deletes a flag by its ID.
//
// @Summary Delete flag
// @Description Deletes one flag by `flag` query parameter (flag code).
// @Tags flags
// @Produce json
// @Security CookieAuth
// @Param flag query string true "Flag code"
// @Success 200 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /delete-flag [delete]
func (h *Handler) HandleDeleteFlag(c fiber.Ctx) error {
	flagID := c.Query("flag")
	if flagID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Error: "Missing flag ID",
		})
	}

	if err := h.store.Queries.DeleteFlagByCode(c.RequestCtx(), flagID); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to delete flag")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Error: "Failed to delete flag",
		})
	}

	return c.JSON(ResponseSuccess{Message: "Flag deleted successfully"})
}

// @Summary List exploits
// @Description Returns all stored exploits.
// @Tags exploits
// @Produce json
// @Security CookieAuth
// @Success 200 {object} map[string]any
// @Failure 500 {object} ResponseError
// @Router /exploits [get]
func (h *Handler) HandleGetExploits(c fiber.Ctx) error {
	exploits, err := h.store.Queries.GetAllExploits(c.RequestCtx())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch exploits")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Error: "Failed to fetch exploits",
		})
	}

	return c.JSON(fiber.Map{
		"exploits": exploits,
		"count":    len(exploits),
	})
}

type ExploitWithContent struct {
	database.Exploit
	Content string `json:"content"`
}

// @Summary Get exploit by name
// @Description Returns exploit(s) with content by name.
// @Tags exploits
// @Produce json
// @Security CookieAuth
// @Param name path string true "Exploit name"
// @Success 200 {array} ExploitWithContent
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /exploit/{name} [get]
func (h *Handler) HandleGetExploit(c fiber.Ctx) error {
	exploitName := c.Params("name")
	if exploitName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Error: "Missing exploit name",
		})
	}

	exploits, err := h.store.Queries.GetExploitsByName(c.RequestCtx(), exploitName)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch exploit by name")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Error: "Failed to fetch exploit",
		})
	}

	if len(exploits) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(ResponseError{
			Error: "Exploit not found",
		})
	}

	exploitsWContent := make([]ExploitWithContent, len(exploits))
	for i, ex := range exploits {
		exploitsWContent[i] = ExploitWithContent{
			Exploit: ex,
			Content: "",
		}

		exploitFile, err := os.OpenFile(path.Join(exploit.ExploitPath, exploits[i].Hash+".py"), os.O_RDONLY, 0o644)
		if err != nil {
			logger.Log.Error().Err(err).Str("hash", exploits[i].Hash).Msg("Failed to open exploit file")
			exploitsWContent[i].Content = "Error loading exploit content"
			continue
		}
		content, err := io.ReadAll(exploitFile)
		if err != nil {
			logger.Log.Error().Err(err).Str("hash", exploits[i].Hash).Msg("Failed to read exploit file")
			exploitsWContent[i].Content = "Error loading exploit content"
			exploitFile.Close()
			continue
		}
		exploitsWContent[i].Content = string(content)
		exploitFile.Close()
	}

	return c.JSON(exploitsWContent)
}

func sanitazeExploit(c fiber.Ctx, fileHeader *multipart.FileHeader) (*multipart.FileHeader, error) {
	const maxExploitSize = 10 << 20
	if fileHeader.Size > maxExploitSize {
		return nil, c.Status(fiber.StatusRequestEntityTooLarge).JSON(ResponseError{Error: "Uploaded file is too large"})
	}

	fileHeader.Filename = filepath.Base(fileHeader.Filename)
	if fileHeader.Filename == "" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(ResponseError{Error: "Invalid file name"})
	}

	return fileHeader, nil
}

// @Summary Upload exploit
// @Description Uploads a new exploit file.
// @Tags exploits
// @Accept multipart/form-data
// @Produce json
// @Security CookieAuth
// @Param file formData file true "Exploit file"
// @Success 200 {object} map[string]any
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 413 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /exploit [post]
func (h *Handler) HandlePostExploit(c fiber.Ctx) error {
	token := c.Cookies("token", "")
	jwtParsed, err := VerifyToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ResponseError{Error: "Invalid token"})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.Log.Warn().Err(err).Msg("No file provided in request")
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Error: "Missing file upload (field name: 'file')"})
	}

	fileHeader, err = sanitazeExploit(c, fileHeader)
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot open file"})
	}
	hash := exploit.GenerateHashFromFile(file)
	ex, err := h.store.Queries.GetExploitByHash(c, hash)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) && !strings.Contains(err.Error(), "no rows in result set") {
			logger.Log.Error().Err(err).Msg("Failed to check existing exploit")
			return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
				Error: "Failed to check existing exploit",
			})
		}
		ex = database.Exploit{}
	}

	if ex.Hash != "" {
		return c.Status(fiber.StatusConflict).JSON(ResponseError{
			Error: "Exploit with the same hash already exists",
		})
	}

	exploitsByName, err := h.store.Queries.GetExploitsByName(c, fileHeader.Filename)
	if err != nil && !errors.Is(err, sql.ErrNoRows) && !strings.Contains(err.Error(), "no rows in result set") {
		logger.Log.Error().Err(err).Msg("Failed to check existing exploit by name")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Error: "Failed to check existing exploit by name",
		})
	}

	var version int64 = 1
	if len(exploitsByName) > 0 {
		version = exploitsByName[0].Version + 1
	}

	err = exploit.CreateExploitFile(c, fileHeader, hash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot save file"})
	}

	exploitS := database.CreateExploitParams{
		Name:       fileHeader.Filename,
		Hash:       hash,
		Username:   jwtParsed.Claims.(jwt.MapClaims)["username"].(string),
		SubmitTime: time.Now().Unix(),
		Version:    version,
	}

	err = h.store.Queries.CreateExploit(c, exploitS)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot save exploit metadata"})
	}

	return c.JSON(fiber.Map{
		"message":      "uploaded successfully",
		"exploit_name": fileHeader.Filename,
		"hash":         hash,
		"version":      version,
	})
}

// @Summary Delete exploit
// @Description Deletes an exploit by ID.
// @Tags exploits
// @Produce json
// @Security CookieAuth
// @Param id path int true "Exploit ID"
// @Success 200 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /exploit/{id} [delete]
func (h *Handler) HandleDeleteExploit(c fiber.Ctx) error {
	exploitID := c.Params("id")
	if exploitID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Error: "Missing exploit ID",
		})
	}

	id, err := strconv.ParseInt(exploitID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Error: "Invalid exploit ID",
		})
	}
	err = h.store.Queries.DeleteExploitByID(c.RequestCtx(), id)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to delete exploit")
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Error: "Failed to delete exploit",
		})
	}

	return c.JSON(ResponseSuccess{Message: "Exploit deleted successfully"})
}
