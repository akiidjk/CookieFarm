package config

import "github.com/ByteTheCookies/cookiefarm-client/internal/models"

const CYCLE_TIME = 15 // seconds
var FAKE_SERVICES = []models.Service{
	{Name: "CCApp", Port: 80},
	{Name: "Ticket", Port: 1337},
	{Name: "Poll", Port: 8080},
	{Name: "COOKIEFLAG", Port: 6969},
}
