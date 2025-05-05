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

// New initializes and returns a new Fiber app instance,
// setting up static file routes, debug middleware, and template engine.
func New() *fiber.App {
	cfg := newConfig(*config.Debug)
	app := fiber.New(cfg)

	// Serve static assets from public folders with compression and caching
	app.Static("/css", "./public/css", fiber.Static{
		Compress:      true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})
	app.Static("/js", "./public/js", fiber.Static{
		Compress:      true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})
	app.Static("/images", "./public/images", fiber.Static{
		Compress:      true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

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
