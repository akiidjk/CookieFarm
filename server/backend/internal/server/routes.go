package server

import (
	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/protocols"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.GetStatus)
	s.App.Get("/stats", s.GetStats)
	s.App.Get("/flags", s.GetFlags)
	s.App.Get("/get-config", s.GetConfig)
	s.App.Get("/health", s.healthHandler)

	s.App.Post("/submit-flags", s.SubmitFlag)

}

func (s *FiberServer) GetConfig(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"config": map[string]interface{}{
			"max_flags": 10,
			"max_users": 100,
		},
	})
}

func (s *FiberServer) SubmitFlag(c *fiber.Ctx) error {
	logger.Debug("SUBMITFLAG | Request received by %s", c.IP())
	body := models.FlagResponse{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	// logger.Debug("Body parsed %v", body)
	s.db.AddFlags(body.Flags)
	flags, err := s.db.GetFlagsCode()
	if err != nil {
		logger.Error("Error %v", err)
	}

	res, err := protocols.Submit(flags)
	if err != nil {
		logger.Error("Error %v", err)
	}

	logger.Debug("Results: %v", res)

	return c.JSON(fiber.Map{
		"message": "Flag submitted successfully",
	})
}

func (s *FiberServer) GetStats(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"stats": map[string]interface{}{
			"total_flags": 0,
			"total_users": 0,
		},
	})
}

func (s *FiberServer) GetFlags(c *fiber.Ctx) error {
	logger.Debug("GETFLAGS | Request received by %s", c.IP())
	flags, err := s.db.GetFlags()
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
