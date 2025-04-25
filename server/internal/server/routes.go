package server

import (
	"context"
	"math"
	"strconv"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
)

func cookieAuthMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" || VerifyToken(token) != nil {
		return c.Redirect("/login")
	}
	return nil
}

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true, // credentials require explicit origins
		MaxAge:           300,
	}))

	publicView := s.App.Group("/")
	publicView.Get("/", s.HandleIndexPage)
	publicView.Get("/dashboard", s.HandleIndexPage)
	publicView.Get("/login", s.HandleLoginPage)
	publicView.Get("/flags/:limit", s.HandlePaginatedFlags)
	publicView.Get("/pagination/:limit", s.HandlePagination)

	publicApi := s.App.Group("/api/v1")
	publicApi.Get("/", s.GetStatus)
	publicApi.Post("/auth/login", NewLimiter(), s.HandleLogin)

	// Aspected Header with: `Authorization: Bearer <token>`
	privateApi := s.App.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey:  config.Secret,
		TokenLookup: "header:Authorization,cookie:token",
	}))
	privateApi.Get("/stats", s.GetStats)
	privateApi.Get("/flags", s.GetAllFlags)
	privateApi.Get("/flags/:limit", s.GetPaginatedFlags)
	privateApi.Get("/config", s.GetConfig)
	privateApi.Get("/health", s.healthHandler)
	privateApi.Post("/submit-flags", s.SubmitFlags)
	privateApi.Post("/submit-flag", s.SubmitFlag)
	privateApi.Post("/config", s.SetConfig)

}

func (s *FiberServer) HandleIndexPage(c *fiber.Ctx) error {

	if err := cookieAuthMiddleware(c); err != nil {
		return err
	}

	limit := c.QueryInt("limit", 100)
	if limit <= 0 {
		limit = 100
	}
	totalFlags, err := s.db.FlagsNumber(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Errore nel recupero dei dati")
	}
	logger.Log.Info().Int("Limit", limit).Msg("Index page request")
	return c.Render("pages/dashboard", fiber.Map{
		"title": "Dashboard",
		"Pagination": struct {
			Limit int
			Pages int
		}{
			Limit: limit,
			Pages: totalFlags / limit,
		},
	}, "layouts/main")
}

func (s *FiberServer) HandlePagination(c *fiber.Ctx) error {
	if err := cookieAuthMiddleware(c); err != nil {
		return err
	}
	limitStr := c.Params("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Parametro limit non valido")
	}

	logger.Log.Debug().Int("limit", limit).Msg("Paginated flags request")

	totalFlags, err := s.db.FlagsNumber(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Errore nel recupero dei dati")
	}

	offset := c.QueryInt("offset", 0)

	totalPages := int(math.Ceil(float64(totalFlags) / float64(limit)))
	current := offset / limit
	pageList := utils.MakePagination(current, totalPages)

	return c.Render("partials/pagination", fiber.Map{
		"Pagination": models.Pagination{
			Limit:    limit,
			Pages:    totalPages,
			Current:  current,
			HasPrev:  current > 0,
			HasNext:  current < totalPages-1,
			PageList: pageList,
		},
	})
}

func (s *FiberServer) HandlePaginatedFlags(c *fiber.Ctx) error {
	if err := cookieAuthMiddleware(c); err != nil {
		return err
	}
	limitStr := c.Params("limit")
	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Parametro limit non valido")
	}

	offsetStr := c.Query("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Parametro offset non valido")
	}

	logger.Log.Debug().Int("offset", offset).Int("limit", limit).Msg("Paginated flags request")

	flags, err := s.db.GetPagedFlags(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Errore nel recupero dei dati")
	}

	logger.Log.Debug().Int("n_flags", len(flags)).Msg("Paginated flags response")

	return c.Render("partials/flags_rows", fiber.Map{
		"Flags": flags,
	})

}

func (s *FiberServer) HandleLoginPage(c *fiber.Ctx) error {
	return c.Render("pages/login", fiber.Map{
		"title": "Login",
	}, "layouts/main")
}

func (s *FiberServer) GetConfig(c *fiber.Ctx) error {
	return c.JSON(config.Current)
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

func (s *FiberServer) GetAllFlags(c *fiber.Ctx) error {
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

func (s *FiberServer) GetPaginatedFlags(c *fiber.Ctx) error {
	limitStr := c.Params("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit parameter",
		})
	}

	offsetStr := c.Query("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid offset parameter",
		})
	}

	flags, err := s.db.GetPagedFlags(limit, offset)
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
