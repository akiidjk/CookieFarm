package api

import (
	"database/sql"
	"logger"
	"models"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"server/config"
	"server/controllers"
	"server/database"
	"server/websockets"

	json "github.com/bytedance/sonic"

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
// Currently returns placeholders for flags and users.
//
// @Summary Get stats
// @Description Returns aggregated server/flags statistics.
// @Tags stats
// @Produce json
// @Security CookieAuth
// @Success 200 {object} map[string]any
// @Failure 500 {object} ResponseError
// @Router /stats [get]
func (*Handler) HandleGetStats(c fiber.Ctx) error {
	n := controllers.NewStatsController()
	return n.GetFlagStats(c)
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
			String: c.Query("search"),
			Valid:  c.Query("search") != "",
		},
		SearchField: sql.NullString{
			String: c.Query("search_field"),
			Valid:  c.Query("search_field") != "",
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
		Search:      c.Query("search", ""),                                    // Value of the search query
		SearchField: c.Query("search_field", "flag_code"),                     // Field to apply the search query to (default: flag_code)
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
	for _, path := range searchPaths {
		if protocols, err := os.ReadDir(path); err == nil {
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

	cfgJSON, err := json.Marshal(h.config.GetFullConfig())
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

// fiber:context-methods migrated
