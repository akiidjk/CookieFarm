package server

import (
	"context"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/database"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/gofiber/fiber/v2"
)

// ---------- GET ----------------

// HandleGetConfig returns the current configuration of the server.
func HandleGetConfig(c *fiber.Ctx) error {
	return c.JSON(config.Current)
}

// HandleGetAllFlags retrieves and returns all the stored flags.
func HandleGetAllFlags(c *fiber.Ctx) error {
	flags, err := database.GetAllFlags()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{
			Error: err.Error(),
		})
	}
	if flags == nil {
		flags = []models.Flag{}
	}
	data := models.ResponseFlags{
		Nflags: len(flags),
		Flags:  flags,
	}
	return c.JSON(data)
}

// HandleGetStats returns statistics about the server state.
// Currently returns placeholders for flags and users.
func HandleGetStats(c *fiber.Ctx) error {
	logger.Log.Debug().Msg("Stats endpoint hit")
	return c.JSON(fiber.Map{
		"stats": map[string]any{
			"total_flags": 0,
			"total_users": 0,
		},
	})
}

// HandleGetPaginatedFlags returns a paginated list of flags based on the limit and offset.
func HandleGetPaginatedFlags(c *fiber.Ctx) error {
	limit, err := c.ParamsInt("limit", config.DefaultLimit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{
			Error: "Invalid limit parameter",
		})
	}
	offset := c.QueryInt("offset", config.DefaultOffeset)

	flags, err := database.GetPagedFlags(uint(limit), uint(offset))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{
			Error: err.Error(),
		})
	}

	if flags == nil {
		flags = []models.Flag{}
	}
	data := models.ResponseFlags{
		Nflags: len(flags),
		Flags:  flags,
	}

	return c.JSON(data)
}

// ---------- POST ----------------

// HandlePostFlags processes a batch of flags submitted in the request.
func HandlePostFlags(c *fiber.Ctx) error {
	payload := new(models.SubmitFlagsRequest)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(models.ResponseError{Error: err.Error()})
	}

	flags := payload.Flags

	if err := database.AddFlags(flags); err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Failed to insert flags into DB")
		return c.Status(fiber.StatusInternalServerError).
			JSON(models.ResponseError{Error: "DB error: " + err.Error()})
	}

	payload.Flags = nil

	return c.JSON(models.ResponseSuccess{
		Message: "Flags submitted successfully",
	})
}

// HandlePostFlag processes a single flag and optionally submits it to an external checker.
func HandlePostFlag(c *fiber.Ctx) error {
	payload := new(models.SubmitFlagRequest)

	if err := c.BodyParser(payload); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid SubmitFlag payload")
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(models.ResponseError{Error: err.Error()})
	}
	f := payload.Flag
	if err := database.AddFlag(f); err != nil {
		logger.Log.Error().Err(err).Msg("DB insert failed in SubmitFlag")
		return c.Status(fiber.StatusInternalServerError).
			JSON(models.ResponseError{Error: "Failed to add flag: " + err.Error()})
	}

	flags := []string{f.FlagCode}

	if config.Current.ConfigServer.HostFlagchecker == "" {
		logger.Log.Warn().Msg("Flagchecker host not configured")
		return c.Status(fiber.StatusServiceUnavailable).JSON(models.ResponseError{
			Error: "Flagchecker host not configured",
		})
	}

	response, err := config.Submit(config.Current.ConfigServer.HostFlagchecker, config.Current.ConfigServer.TeamToken, flags)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to submit flag to external checker")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{
			Error:   "Failed to submit flag",
			Details: err.Error(),
		})
	}

	UpdateFlags(response)

	return c.JSON(models.ResponseSuccess{
		Message: "Flag submitted successfully",
	})
}

// HandlePostConfig updates the server configuration and restarts the flag processing loop.
func HandlePostConfig(c *fiber.Ctx) error {
	var configPayload struct {
		Config models.Config `json:"config"`
	}
	if err := c.BodyParser(&configPayload); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to parse config payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ResponseError{
			Error:   "Failed to parse config payload",
			Details: err.Error(),
		})
	}

	config.Current = configPayload.Config

	if shutdownCancel != nil {
		shutdownCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	shutdownCancel = cancel

	go StartFlagProcessingLoop(ctx)

	return c.JSON(models.ResponseSuccess{
		Message: "Configuration updated successfully",
	})
}

func HandleDeleteFlag(c *fiber.Ctx) error {
	flagID := c.Query("flag")
	if flagID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{
			Error: "Flag ID is required",
		})
	}

	if err := database.DeleteFlag(flagID); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to delete flag")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{
			Error: "Failed to delete flag",
		})
	}

	return c.JSON(models.ResponseSuccess{
		Message: "Flag deleted successfully",
	})
}
