package server

import (
	"context"
	"strings"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/database"
	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type FiberServer struct {
	*fiber.App
	shutdownCancel context.CancelFunc
	db             database.Service
}

func DevConfig() fiber.Config {
	return fiber.Config{
		AppName:               "CookieFarm Backend (Dev)",
		DisableStartupMessage: false,
		Prefork:               false,
		CaseSensitive:         false,
		StrictRouting:         false,
		ServerHeader:          "Fiber",
		EnablePrintRoutes:     true,
		Views:                 html.New(ui.GetPathView(), ".html"),
	}
}

func ProdConfig() fiber.Config {
	return fiber.Config{
		AppName:               "CookieFarm Backend",
		DisableStartupMessage: true,
		Prefork:               false,
		CaseSensitive:         true,
		StrictRouting:         true,
		ServerHeader:          "",
		EnablePrintRoutes:     false,
		Views:                 html.NewFileSystem(ui.ViewsFS(), ".html"),
	}
}

func New() *FiberServer {
	var cfg fiber.Config
	if *config.Debug {
		cfg = DevConfig()
	} else {
		cfg = ProdConfig()
	}

	app := fiber.New(cfg)

	app.Static("/css", "./public/css")
	app.Static("/js", "./public/js")
	app.Static("/static", "./public")

	if *config.Debug {
		app.Use(func(c *fiber.Ctx) error {
			path := c.Path()
			if strings.HasPrefix(path, "/css") || strings.HasPrefix(path, "/js") || strings.HasPrefix(path, "/static") {
				logger.Log.Debug().Str("path", path).Msg("Serving static asset")
			}
			return c.Next()
		})
	}

	db := database.New()
	logger.Log.Info().Msg("Database initialized")

	return &FiberServer{
		App: app,
		db:  db,
	}
}

func (s *FiberServer) DB() database.Service {
	return s.db
}
