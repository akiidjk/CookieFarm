package server

import (
	"context"
	"strings"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/database"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/ui"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
	db             database.Service
	shutdownCancel context.CancelFunc
}

func newConfig(debug bool) fiber.Config {
	views := ui.InitTemplateEngine(!debug)
	common := fiber.Config{
		Views: views,
	}

	if debug {
		common.AppName = "CookieFarm Server (Dev)"
		common.DisableStartupMessage = false
		common.CaseSensitive = false
		common.StrictRouting = false
		common.EnablePrintRoutes = true
	} else {
		common.AppName = "CookieFarm Server"
		common.DisableStartupMessage = true
		common.CaseSensitive = true
		common.StrictRouting = true
		common.EnablePrintRoutes = false
	}

	common.Prefork = false
	common.ServerHeader = "Fiber"

	return common
}

func New() *FiberServer {
	cfg := newConfig(*config.Debug)
	app := fiber.New(cfg)

	app.Static("/css", "./public/css")
	app.Static("/js", "./public/js")
	app.Static("/images", "./public/images")

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
