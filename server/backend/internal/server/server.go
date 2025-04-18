package server

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/database"
)

type FiberServer struct {
	*fiber.App
	loopCancel context.CancelFunc
	db         database.Service
}

func DevConfig() fiber.Config {
	return fiber.Config{
		AppName:               "Backend (Development)",
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
		AppName:               "Backend",
		DisableStartupMessage: true,
		Prefork:               false, // Multiprocess
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

	server := &FiberServer{
		App: app,
		db:  database.New(),
	}

	return server
}

func (s *FiberServer) DB() database.Service {
	return s.db
}
