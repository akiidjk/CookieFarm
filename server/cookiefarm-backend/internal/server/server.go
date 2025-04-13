package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ByteTheCookies/cookiefarm-backend/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "github.com/ByteTheCookies/cookiefarm-backend",
			AppName:      "github.com/ByteTheCookies/cookiefarm-backend",
		}),

		db: database.New(),
	}

	return server
}
