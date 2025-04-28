package server

import (
	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
)

func RegisterRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     utils.GetEnv("ALLOW_ORIGINS", "http://localhost:8080"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true, // credentials require explicit origins
		MaxAge:           300,
	}))

	// ---------- VIEW ----------

	View := app.Group("/")
	View.Get("/", HandleIndexPage)
	View.Get("/dashboard", HandleIndexPage)
	View.Get("/login", HandleLoginPage)
	View.Get("/flags/:limit", HandlePartialsFlags)
	View.Get("/pagination/:limit", HandlePartialsPagination)

	// ---------- API ----------

	publicApi := app.Group("/api/v1")
	publicApi.Get("/", GetStatus)
	publicApi.Post("/auth/login", NewLimiter(), HandleLogin)

	// Aspected Header with: `Authorization: Bearer <token>`
	privateApi := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey:  config.Secret,
		TokenLookup: "header:Authorization,cookie:token",
	}))
	privateApi.Get("/stats", HandleGetStats)
	privateApi.Get("/flags", HandleGetAllFlags)
	privateApi.Get("/flags/:limit", HandleGetPaginatedFlags)
	privateApi.Get("/config", HandleGetConfig)
	privateApi.Get("/health", HealthHandler)
	privateApi.Post("/submit-flags", HandlePostFlags)
	privateApi.Post("/submit-flag", HandlePostFlag)
	privateApi.Post("/config", HandlePostConfig)

}

func GetStatus(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "The cookie is up!!",
	}
	return c.JSON(resp)
}

func HealthHandler(c *fiber.Ctx) error {
	return c.JSON(GetStatus(c))
}
