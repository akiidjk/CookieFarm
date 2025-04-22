package server

import (
	"context"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	public := s.App.Group("/api/v1")
	public.Get("/", s.GetStatus)
	public.Post("/auth/login", NewLimiter(), s.HandleLogin)
	public.Post("/auth/verify", NewLimiter(), s.HandleVerify)

	// Aspected Header with: `Authorization: Bearer <token>`
	private := s.App.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: config.Secret,
	}))
	private.Get("/stats", s.GetStats)
	private.Get("/flags", s.GetFlags)
	private.Get("/config", s.GetConfig)
	private.Get("/health", s.healthHandler)
	private.Post("/submit-flags", s.SubmitFlags)
	private.Post("/submit-flag", s.SubmitFlag)
	private.Post("/config", s.SetConfig)

}

func (s *FiberServer) GetConfig(c *fiber.Ctx) error {
	return c.JSON(config.Current)
}

func (s *FiberServer) HandleVerify(c *fiber.Ctx) error {
	var verifyPayload struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&verifyPayload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	if err := VerifyToken(verifyPayload.Token); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Token verified successfully",
	})
}

func (s *FiberServer) SetConfig(c *fiber.Ctx) error {
	var configPayload struct {
		Config models.Config `json:"config"`
	}
	if err := c.BodyParser(&configPayload); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to parse config payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
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

	return c.JSON(fiber.Map{
		"message": "Configuration updated successfully",
	})
}

func (s *FiberServer) SubmitFlag(c *fiber.Ctx) error {
	body := map[string]models.Flag{"flag": models.Flag{}}
	if err := c.BodyParser(&body); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid SubmitFlag payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := s.db.AddFlag(body["flag"]); err != nil {
		logger.Log.Error().Err(err).Msg("DB insert failed in SubmitFlag")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add flag: " + err.Error(),
		})
	}

	flags := []string{body["flag"].FlagCode}
	response, err := config.Submit(config.Current.ConfigServer.HostFlagchecker, config.Current.ConfigServer.TeamToken, flags)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to submit flag to external checker")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to submit flag",
			"details": err.Error(),
		})
	}

	logger.Log.Info().Strs("submitted_flags", flags).Msg("Flag submitted successfully")
	s.UpdateFlags(response)

	return c.JSON(fiber.Map{
		"message": "Flag submitted successfully",
	})
}

func (s *FiberServer) GetStats(c *fiber.Ctx) error {
	logger.Log.Debug().Msg("Stats endpoint hit")
	return c.JSON(fiber.Map{
		"stats": map[string]interface{}{
			"total_flags": 0,
			"total_users": 0,
		},
	})
}

func (s *FiberServer) SubmitFlags(c *fiber.Ctx) error {
	body := map[string][]models.Flag{"flags": {}}
	if err := c.BodyParser(&body); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid SubmitFlags payload")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := s.db.AddFlags(body["flags"]); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to insert flags into DB")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "DB error: " + err.Error(),
		})
	}

	logger.Log.Info().Int("count", len(body["flags"])).Msg("Flags batch submitted")
	return c.JSON(fiber.Map{
		"message": "Flags submitted successfully",
	})
}

func (s *FiberServer) GetFlags(c *fiber.Ctx) error {
	flags, err := s.db.GetAllFlags()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if flags == nil {
		flags = []models.Flag{}
	}
	return c.JSON(fiber.Map{
		"n_flags": len(flags),
		"flags":   flags,
	})
}

func (s *FiberServer) GetStatus(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "The cookie is up!!",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
