package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/server"
	"github.com/ByteTheCookies/backend/internal/utils"

	flogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.Debug = flag.Bool("debug", false, "Abilita log livello debug")
	flag.Parse()

	if *config.Debug {
		logger.Setup("debug")
	} else {
		logger.Setup("info")
	}
	defer logger.Close()

	server.InitSecret()
	logger.Log.Debug().Str("password_plain", config.Password).Msg("Password before hash")

	var err error
	config.Password, err = server.HashPassword(config.Password)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to hash password")
	}
	logger.Log.Debug().Str("password_hash", config.Password).Msg("Password after hash")

	app := server.New()

	if *config.Debug {
		app.Use(flogger.New(flogger.Config{
			Format:     "[${time}] ${ip} - ${method} ${path} - ${status}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "Local",
		}))
	}

	app.RegisterFiberRoutes()

	done := make(chan bool, 1)
	go gracefulShutdown(app, done)

	port, _ := strconv.Atoi(utils.GetEnv("PORT", "8080"))
	go func() {
		addr := fmt.Sprintf(":%d", port)
		logger.Log.Info().Str("addr", addr).Msg("Starting HTTP server")
		if err := app.Listen(addr); err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	<-done
	logger.Log.Info().Msg("Graceful shutdown complete.")
}

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Log.Warn().Msg("Shutting down gracefully...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctxTimeout); err != nil {
		logger.Log.Error().Err(err).Msg("Forced server shutdown")
	}

	done <- true
}
