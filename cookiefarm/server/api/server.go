// Package api initializes and configures the HTTP server for CookieFarm,
// including routing, static file serving, and debug settings.
package api

import (
	"fmt"
	"logger"
	"os"
	"strings"
	"time"

	"server/config"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/static"
)

const (
	frontendDistDir   = "./server/frontend/dist"
	frontendAssetsDir = "./server/frontend/dist/assets"
	frontendIndexPath = "./server/frontend/dist/index.html"
)

// newConfig returns a configured Fiber config struct.
// It adapts settings depending on the debug mode (e.g. logging, strict routing).
func newConfig(debug bool) fiber.Config {
	cfg := fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	}

	if debug {
		cfg.AppName = "CookieFarm Server (Dev)"
		cfg.CaseSensitive = false
		cfg.StrictRouting = false
	} else {
		cfg.AppName = "CookieFarm Server"
		cfg.CaseSensitive = true
		cfg.StrictRouting = true
	}

	cfg.ServerHeader = "Fiber" // Custom server header.

	return cfg
}

func PrepareStatic(app *fiber.App) error {
	type staticRoute struct {
		route string
		dir   string
	}

	routes := []staticRoute{
		{"/images", "./server/public/images"},
		{"/assets", frontendAssetsDir},
	}

	var staticCfg static.Config
	if config.Cache {
		staticCfg = static.Config{
			Compress:      true,
			CacheDuration: 10 * time.Second,
			MaxAge:        3600,
		}
	} else {
		staticCfg = static.Config{
			Compress: true,
		}
	}

	for _, r := range routes {
		app.Get(r.route+"/*", static.New(r.dir, staticCfg))
	}

	return nil
}

func ServeFrontendIndex(c fiber.Ctx) error {
	if _, err := os.Stat(frontendIndexPath); err != nil {
		return c.Status(fiber.StatusNotFound).SendString("frontend build not found")
	}

	return c.SendFile(frontendIndexPath)
}

// NewApp initializes and returns a new Fiber app instance,
// setting up static file routes, debug middleware, and template engine.
func NewApp() (*fiber.App, error) {
	staticPrefixes := []string{"/css", "/js", "/images", "/assets", "/static"}
	cfg := newConfig(config.Debug)
	app := fiber.New(cfg)

	if err := PrepareStatic(app); err != nil {
		return nil, fmt.Errorf("prepare static: %w", err)
	}

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Server().ReadTimeout = 10 * time.Second
	app.Server().WriteTimeout = 10 * time.Second
	app.Server().IdleTimeout = 60 * time.Second

	// Log static file requests in debug mode
	if config.Debug {
		app.Use(func(c fiber.Ctx) error {
			for _, prefix := range staticPrefixes {
				if strings.HasPrefix(c.Path(), prefix) {
					logger.Log.Debug().Str("path", c.Path()).Msg("Serving static asset")
					break
				}
			}
			return c.Next()
		})
	}

	return app, nil
}
