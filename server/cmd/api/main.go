package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/server"
	"github.com/ByteTheCookies/cookieserver/internal/utils"

	flogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.Debug = flag.Bool("debug", false, "Enable debug-level logging")
	flag.Parse()

	level := "info"
	if *config.Debug {
		level = "debug"
	}
	logger.Setup(level)
	defer logger.Close()

	server.InitSecret()
	logger.Log.Debug().Str("plain", config.Password).Msg("Plain password before hashing")

	hashed, err := server.HashPassword(config.Password)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Password hashing failed")
	}
	config.Password = hashed
	logger.Log.Debug().Str("hashed", config.Password).Msg("Password after hashing")

	app := server.New()
	if *config.Debug {
		app.Use(flogger.New(flogger.Config{
			Format:     "[${time}] ${ip} - ${method} ${path} - ${status}\n",
			TimeFormat: time.RFC3339,
			TimeZone:   "Local",
		}))
	}
	app.RegisterRoutes()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := utils.GetEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	go func() {
		logger.Log.Info().Str("addr", addr).Msg("HTTP server starting")
		if err := app.Listen(addr); err != nil {
			logger.Log.Fatal().Err(err).Msg("Server listen error")
		}
	}()

	<-ctx.Done()
	logger.Log.Warn().Msg("Shutdown signal received, terminating...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("Error during shutdown, forcing exit")
	}

	logger.Log.Info().Msg("Server stopped gracefully")
}
