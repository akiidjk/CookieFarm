// Package config provides functions to manage the CookieFarm client configuration globally.
package config

import "github.com/ByteTheCookies/cookieclient/internal/models"

var Current models.Config // Global current configuration
var Token string          // Global token for the authentication
