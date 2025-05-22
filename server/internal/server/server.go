// Package server initializes and configures the HTTP server for CookieFarm,
// including routing, static file serving, and debug settings.
package server

import (
	"context"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/ui"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

// shutdownCancel is used to gracefully shut down the server from external signals (if implemented).
var shutdownCancel context.CancelFunc

// newConfig returns a configured Fiber config struct.
// It adapts settings depending on the debug mode (e.g. logging, strict routing).
func newConfig(debug bool) fiber.Config {
	views := ui.InitTemplateEngine(!debug)

	cfg := fiber.Config{
		Views:       views,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	}

	if debug {
		cfg.AppName = "CookieFarm Server (Dev)"
		cfg.DisableStartupMessage = false
		cfg.CaseSensitive = false
		cfg.StrictRouting = false
		cfg.EnablePrintRoutes = true
	} else {
		cfg.AppName = "CookieFarm Server"
		cfg.DisableStartupMessage = false
		cfg.CaseSensitive = true
		cfg.StrictRouting = true
		cfg.EnablePrintRoutes = false
	}

	cfg.Prefork = false        // Disable prefork mode (multi-process); not needed here.
	cfg.ServerHeader = "Fiber" // Custom server header.

	return cfg
}

func PrepareStatic(app *fiber.App) error {
	// Serve static assets from public folders with compression and caching

	type staticRoute struct {
		route string
		dir   string
	}

	routes := []staticRoute{
		{"/css", "./public/css"},
		{"/js", "./public/js"},
		{"/images", "./public/images"},
	}

	var staticCfg fiber.Static
	if config.Cache {
		staticCfg = fiber.Static{
			Compress:      true,
			CacheDuration: 10 * time.Second,
			MaxAge:        3600,
		}
	} else {
		staticCfg = fiber.Static{
			Compress: true,
		}
	}

	for _, r := range routes {
		app.Static(r.route, r.dir, staticCfg)
	}

	return nil
}

// New initializes and returns a new Fiber app instance,
// setting up static file routes, debug middleware, and template engine.
func New() *fiber.App {
	cfg := newConfig(*config.Debug)
	app := fiber.New(cfg)

	// Serve static assets from public folders with compression and caching
	if err := PrepareStatic(app); err != nil {
		logger.Log.Error().Err(err).Msg("Error preparing static assets")
	}

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Server().ReadTimeout = 10 * time.Second
	app.Server().WriteTimeout = 10 * time.Second

	// Log static file requests in debug mode
	if *config.Debug {
		app.Use(func(c *fiber.Ctx) error {
			path := c.Path()
			if strings.HasPrefix(path, "/css") || strings.HasPrefix(path, "/js") || strings.HasPrefix(path, "/static") {
				logger.Log.Debug().Str("path", path).Msg("Serving static asset")
			}
			return c.Next()
		})
	}

	return app
}
