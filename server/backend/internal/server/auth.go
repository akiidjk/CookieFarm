package server

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func InitSecret() {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		panic(fmt.Sprintf("failed to generate secret: %v", err))
	}
	config.Secret = secret
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CreateJWT() (string, int64, error) {
	exp := time.Now().Add(time.Hour * 24).Unix()
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

func (s *FiberServer) Login(c *fiber.Ctx) error {
	req := new(models.SigninRequest)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(config.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, exp, err := CreateJWT()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": token, "exp": exp})
}
