// Package main is the entry point for the API server.
package main

import (
	"context"
	_ "embed"
	"fmt"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/database"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/server"
	"github.com/ByteTheCookies/cookieserver/internal/utils"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/pflag"
)

//go:embed banner.txt
var banner string

func init() {
	fmt.Println(banner)

	config.Debug = pflag.BoolP("debug", "d", false, "Enable debug-level logging")
	config.ConfigPath = *pflag.StringP("config", "c", "", "Path to the configuration file")
	config.Password = *pflag.StringP("password", "p", "password", "Password for authentication")
	config.ServerPort = *pflag.StringP("port", "P", "8080", "Port for server")
}

// The main function initializes configuration, sets up logging, connects to the database,
// configures the Fiber HTTP server, and handles graceful shutdown on system signals.
func main() {
	pflag.Parse()

	level := "info"
	if *config.Debug {
		level = "debug"
	}
	logger.Setup(level)
	defer logger.Close()

	if config.ConfigPath != "" {
		logger.Log.Info().Msg("Using file config...")
		err := server.LoadConfig(config.ConfigPath)
		if err != nil {
			logger.Log.Warn().Err(err).Msg("Config file not found or corrupted using web config")
		}
	} else {
		logger.Log.Info().Msg("Using web config...")
	}

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

	server.RegisterRoutes(app)

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

	database.DB = database.New()
	logger.Log.Info().Msg("Database initialized")
	defer database.Close()

	<-ctx.Done()
	logger.Log.Warn().Msg("Shutdown signal received, terminating...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("Error during shutdown, forcing exit")
	}

	logger.Log.Info().Msg("Server stopped gracefully")
}
