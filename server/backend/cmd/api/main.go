package main

import (
	"context"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/server"
	"github.com/ByteTheCookies/backend/internal/utils"

	flogger "github.com/gofiber/fiber/v2/middleware/logger"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	logger.Warning("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown with error: %v", err)
	}

	logger.Warning("Server exiting")

	done <- true
}

func init() {
	logger.SetLevel(0)
}

func main() {
	server := server.New()

	if utils.GetEnv("IsDevelopment", "true") == "true" {
		server.Use(flogger.New(flogger.Config{
			Format:     "[${time}] ${ip} - ${method} ${path} - ${status}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "Local",
		}))
	}

	server.RegisterFiberRoutes()

	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(utils.GetEnv("PORT", "8080"))
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	go gracefulShutdown(server, done)

	<-done
	logger.Info("Graceful shutdown complete.")
}
