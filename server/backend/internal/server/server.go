package server

import (
	"context"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/database"
	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/gofiber/fiber/v2"
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
	}
}

func New() *FiberServer {
	var app *fiber.App
	if *config.Debug {
		app = fiber.New(DevConfig())
	} else {
		app = fiber.New(ProdConfig())
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
