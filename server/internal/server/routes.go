package server

import (
	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/database"
	"github.com/ByteTheCookies/cookieserver/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
)

// RegisterRoutes configures all routes and middlewares of the Fiber app,
// including CORS policies, public and protected API endpoints, and view rendering routes.
func RegisterRoutes(app *fiber.App) {
	// Enable CORS with dynamic origins from environment variable.
	// Useful for allowing access from web dashboards or dev clients.
	app.Use(cors.New(cors.Config{
		AllowOrigins:     utils.GetEnv("ALLOW_ORIGINS", "http://localhost:8080"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ------------------ VIEW ROUTES ------------------

	// These endpoints render HTML templates (SSR views)
	View := app.Group("/")
	View.Get("/", HandleIndexPage)
	View.Get("/dashboard", HandleIndexPage)
	View.Get("/login", HandleLoginPage)
	View.Get("/flags/:limit", HandlePartialsFlags)
	View.Get("/pagination/:limit", HandlePartialsPagination)

	// ------------------ PUBLIC API ------------------

	// Public endpoints that do not require authentication
	publicApi := app.Group("/api/v1")
	publicApi.Get("/", GetStatus) // Simple status check
	publicApi.Post("/auth/login", NewLimiter(), HandleLogin)

	// ------------------ PRIVATE API ------------------

	// Protected routes, accessible only with a valid JWT token.
	// The token can be passed via Authorization header or cookie.
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

// GetStatus is a simple public endpoint used to check if the server is online.
func GetStatus(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "The cookie is up!!",
	}
	return c.JSON(resp)
}

// HealthHandler returns a basic health check result from the database layer.
func HealthHandler(c *fiber.Ctx) error {
	return c.JSON(database.Health())
}
