package config

import (
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/internal/utils"
)

var Current models.Config = models.Config{
	Server: models.ConfigServer{},
	Client: models.ConfigClient{},
}
var Debug *bool
var Secret = make([]byte, 32)
var Password = utils.GetEnv("PASSWORD", "password")
var Submit func(string, string, []string) ([]models.ResponseProtocol, error)
