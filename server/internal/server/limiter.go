package server

import (
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func NewLimiter() fiber.Handler {
	if *config.Debug {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			logger.Log.Warn().
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Msg("Rate limit exceeded on login")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many login attempts. Please try again later.",
			})
		},
	})
}
