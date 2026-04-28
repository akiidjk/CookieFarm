package api

import (
	"logger"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/extractors"

	"server/ckp"
	"server/config"
	"server/core"
	"server/database"
	"server/websockets"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type Handler struct {
	store       *database.Store
	runner      *core.Runner
	config      *config.ConfigManager
	connections *ckp.Connections
}

func NewHandler(s *database.Store, r *core.Runner, c *config.ConfigManager, conns *ckp.Connections) *Handler {
	return &Handler{store: s, runner: r, config: c, connections: conns}
}

// RegisterRoutes configures all routes and middlewares of the Fiber app,
// including CORS policies, public and protected API endpoints, and view rendering routes.
func (h *Handler) RegisterRoutes(app *fiber.App) {
	if len(config.Secret) == 0 {
		logger.Log.Fatal().Msg("JWT secret not configured")
	}

	// Enable CORS with dynamic origins from environment variable.
	// Useful for allowing access from web dashboards or dev clients.
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(config.GetEnv("ALLOW_ORIGINS", "http://localhost:8080,http://localhost:3000,http://localhost:5173"), ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ------------------ PUBLIC API ------------------

	// Public endpoints that do not require authentication
	publicAPI := app.Group("/api/v1")
	publicAPI.Get("/", GetStatus) // Simple status check
	publicAPI.Post("/auth/login", NewLimiter(), HandleLogin)
	publicAPI.Post("/auth/logout", HandleLogout)
	publicAPI.Get("/auth/verify", HandleVerify)
	publicAPI.Get("/protocols", h.HandleGetProtocols)
	publicAPI.Get("/swagger/doc.json", HandleSwaggerDoc)
	publicAPI.Get("/swagger", HandleSwaggerUI)

	// ------------------ PRIVATE API ------------------

	// Protected routes, accessible only with a valid JWT token.
	// The token can be passed via Authorization header or cookie.
	privateAPI := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: config.Secret},
		Extractor: extractors.Chain(
			extractors.FromAuthHeader("Bearer"),
			extractors.FromCookie("token"),
		),
	}))
	privateAPI.Get("/stats", h.HandleGetStats)
	privateAPI.Get("/flags", h.HandleGetAllFlags)
	privateAPI.Get("/flags/:limit", h.HandleGetPaginatedFlags)
	privateAPI.Get("/config", h.HandleGetConfig)
	privateAPI.Get("/config/full", h.HandleGetFullConfig)
	privateAPI.Post("/config", h.HandlePostConfig)
	privateAPI.Post("/submit-flags", h.HandlePostFlags)
	privateAPI.Post("/submit-flag", h.HandlePostFlag)
	privateAPI.Post("/submit-flags-standalone", h.HandlePostFlagsStandalone)
	privateAPI.Delete("/delete-flag", h.HandleDeleteFlag)

	// exploits endpoints
	privateAPI.Get("/exploits", h.HandleGetExploits)
	privateAPI.Get("/exploit/:name", h.HandleGetExploit)
	privateAPI.Post("/exploit/upload", h.HandlePostExploit)
	privateAPI.Delete("/exploit/:id", h.HandleDeleteExploit)

	websockets.GlobalManager = websockets.NewManager()
	app.Use("/ws",
		websockets.CookieAuthMiddleware,
		websockets.WebSocketUpgrade,
	)
	app.Get("/ws", websockets.GlobalManager.ServeWS())

	app.Get("/*", func(c fiber.Ctx) error {
		path := c.Path()
		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/ws") ||
			strings.HasPrefix(path, "/assets/") ||
			strings.HasPrefix(path, "/css/") ||
			strings.HasPrefix(path, "/js/") ||
			strings.HasPrefix(path, "/images/") {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return ServeFrontendIndex(c)
	})
}

// GetStatus is a simple public endpoint used to check if the server is online.
//
// @Summary API status
// @Description Returns a simple status response to confirm the server is online.
// @Tags system
// @Success 200 {object} map[string]string
// @Router / [get]
func GetStatus(c fiber.Ctx) error {
	resp := fiber.Map{
		"message": "The cookie is up!!",
		"time":    time.Now().Format(time.RFC3339),
	}
	return c.JSON(resp)
}
