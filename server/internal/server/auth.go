package server

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func InitSecret() {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		panic(fmt.Sprintf("failed to generate secret: %v", err))
	}
	config.Secret = secret
	logger.Log.Info().Msg("JWT secret generated")
}

func VerifyToken(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return config.Secret, nil
	})
	return err
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CreateJWT() (string, int64, error) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "cookie",
		"exp":      exp,
	})

	tokenString, err := token.SignedString(config.Secret)
	if err != nil {
		return "", 0, err
	}
	return tokenString, exp, nil
}

func (s *FiberServer) HandleLogin(c *fiber.Ctx) error {
	req := new(models.SigninRequest)
	if err := c.BodyParser(req); err != nil {
		logger.Log.Warn().Err(err).Msg("Invalid login payload")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Password == "" {
		logger.Log.Warn().Msg("Missing password in login")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(config.Password), []byte(req.Password)); err != nil {
		logger.Log.Warn().Msg("Login failed: invalid password")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, _, err := CreateJWT()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to generate JWT")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "JWT generation error",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   60 * 60 * 24, // 1 giorno
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func CookieAuthMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" || VerifyToken(token) != nil {
		return c.Redirect("/login")
	}
	return nil
}
