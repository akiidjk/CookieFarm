package server

import (
	"context"
	"sync"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/gofiber/fiber/v2"
)

// ---------- GET ----------------

func (s *FiberServer) HandleGetConfig(c *fiber.Ctx) error {
	return c.JSON(config.Current)
}

func (s *FiberServer) HandleGetAllFlags(c *fiber.Ctx) error {
	flags, err := s.db.GetAllFlags()
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

func (s *FiberServer) HandleGetStats(c *fiber.Ctx) error {
	logger.Log.Debug().Msg("Stats endpoint hit")
	return c.JSON(fiber.Map{
		"stats": map[string]interface{}{
			"total_flags": 0,
			"total_users": 0,
		},
	})
}

func (s *FiberServer) HandleGetPaginatedFlags(c *fiber.Ctx) error {
	limit, err := c.ParamsInt("limit", config.LIMIT)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{
			Error: "Invalid limit parameter",
		})
	}
	offset := c.QueryInt("offset", config.OFFSET)

	flags, err := s.db.GetPagedFlags(limit, offset)
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

func (s *FiberServer) HandlePostFlags(c *fiber.Ctx) error {
	payload := new(models.SubmitFlagsRequest)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(models.ResponseError{Error: err.Error()})
	}

	flags := payload.Flags
	if err := s.db.AddFlags(flags); err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Failed to insert flags into DB")
		return c.Status(fiber.StatusInternalServerError).
			JSON(models.ResponseError{Error: "DB error: " + err.Error()})
	}

	payload.Flags = nil
	flags = nil

	logger.Log.Info().
		Int("count", len(flags)).
		Msg("Flags batch submitted")

	return c.JSON(models.ResponseSuccess{
		Message: "Flags submitted successfully",
	})
}

var submitFlagPool = sync.Pool{
	New: func() any {
		return new(models.SubmitFlagRequest)
	},
}

func (s *FiberServer) HandlePostFlag(c *fiber.Ctx) error {
	payload := new(models.SubmitFlagRequest)

	if err := c.BodyParser(payload); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid SubmitFlag payload")
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(models.ResponseError{Error: err.Error()})
	}
	f := payload.Flag

	if err := s.db.AddFlag(f); err != nil {
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

	logger.Log.Info().Strs("submitted_flags", flags).Msg("Flag submitted successfully")
	s.UpdateFlags(response)

	return c.JSON(models.ResponseSuccess{
		Message: "Flag submitted successfully",
	})
}

func (s *FiberServer) HandlePostConfig(c *fiber.Ctx) error {
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
	logger.Log.Info().Msg("Configuration updated via API")

	if s.shutdownCancel != nil {
		s.shutdownCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdownCancel = cancel

	go s.StartFlagProcessingLoop(ctx)

	return c.JSON(models.ResponseSuccess{
		Message: "Configuration updated successfully",
	})
}
