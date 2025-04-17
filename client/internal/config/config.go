package config

import (
	"strconv"

	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
)

var CYCLE_TIME, _ = strconv.Atoi(utils.GetEnv("CYCLE_TIME", "15")) // seconds
var FAKE_SERVICES = []models.Service{
	{Name: "CCApp", Port: 80},
	{Name: "Ticket", Port: 1337},
	{Name: "Poll", Port: 8080},
	{Name: "COOKIEFLAG", Port: 6969},
}
var BASE_URL_SERVER = utils.GetEnv("BASE_URL_SERVER", "http://localhost:8080")

var Current models.Config
var Token string
