// Package api initializes and configures the HTTP server for CookieFarm,
// including routing, static file serving, and debug settings.
package api

import (
	"fmt"
	"io/fs"
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
	frontendFS := os.DirFS("./server/frontend/dist")
	publicFS := os.DirFS("./server/public")

	assetsFS, err := fs.Sub(frontendFS, "assets")
	if err != nil {
		return err
	}

	imagesFS, err := fs.Sub(publicFS, "images")
	if err != nil {
		return err
	}

	app.Use("/assets", static.New("", static.Config{
		FS:            assetsFS,
		Compress:      true,
		MaxAge:        31536000,
		CacheDuration: 10 * time.Second,
	}))

	app.Use("/images", static.New("", static.Config{
		FS:            imagesFS,
		Compress:      true,
		MaxAge:        3600,
		CacheDuration: 10 * time.Second,
	}))

	app.Use("/", static.New("", static.Config{
		FS:            frontendFS,
		IndexNames:    []string{"index.html"},
		Compress:      true,
		MaxAge:        0,
		CacheDuration: 10 * time.Second,
	}))

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
