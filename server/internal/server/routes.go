package server

import (
	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
)

func (s *FiberServer) RegisterRoutes() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     utils.GetEnv("ALLOW_ORIGINS", "http://localhost:8080"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true, // credentials require explicit origins
		MaxAge:           300,
	}))

	// ---------- VIEW ----------

	View := s.App.Group("/")
	View.Get("/", s.HandleIndexPage)
	View.Get("/dashboard", s.HandleIndexPage)
	View.Get("/login", s.HandleLoginPage)
	View.Get("/flags/:limit", s.HandlePartialsFlags)
	View.Get("/pagination/:limit", s.HandlePartialsPagination)

	// ---------- API ----------

	publicApi := s.App.Group("/api/v1")
	publicApi.Get("/", s.GetStatus)
	publicApi.Post("/auth/login", NewLimiter(), s.HandleLogin)

	// Aspected Header with: `Authorization: Bearer <token>`
	privateApi := s.App.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey:  config.Secret,
		TokenLookup: "header:Authorization,cookie:token",
	}))
	privateApi.Get("/stats", s.HandleGetStats)
	privateApi.Get("/flags", s.HandleGetAllFlags)
	privateApi.Get("/flags/:limit", s.HandleGetPaginatedFlags)
	privateApi.Get("/config", s.HandleGetConfig)
	privateApi.Get("/health", s.HealthHandler)
	privateApi.Post("/submit-flags", s.HandlePostFlags)
	privateApi.Post("/submit-flag", s.HandlePostFlag)
	privateApi.Post("/config", s.HandlePostConfig)

}

func (s *FiberServer) GetStatus(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "The cookie is up!!",
	}
	return c.JSON(resp)
}

func (s *FiberServer) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
