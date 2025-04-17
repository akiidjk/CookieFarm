package config

import (
	"strconv"

	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/internal/utils"
)

var MAX_FLAG_BATCH_SIZE, _ = strconv.Atoi(utils.GetEnv("MAX_FLAG_BATCH_SIZE", "1000"))
var HOST = utils.GetEnv("HOST", "localhost:3000")
var TEAM_TOKEN = utils.GetEnv("TEAM_TOKEN", "4242424242424242")

var Current models.Config = models.Config{
	Server: models.ConfigServer{},
	Client: models.ConfigClient{},
}

var Session models.Session = models.Session{
	Password: "",
}

var Secret = make([]byte, 32)
var Password = utils.GetEnv("PASSWORD", "password")
