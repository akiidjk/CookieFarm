package config

import (
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/internal/utils"
)

var Current models.Config = models.Config{
	Configured:   false,
	ConfigClient: models.ConfigClient{},
	ConfigServer: models.ConfigServer{},
}
var Debug *bool
var Secret = make([]byte, 32)
var Password = utils.GetEnv("PASSWORD", "password")
var ServerPort = utils.GetEnv("BACKEND_PORT", "8080")
var Submit func(string, string, []string) ([]models.ResponseProtocol, error)
