// Package config for configuration management
package config

import (
	"github.com/ByteTheCookies/cookieserver/internal/models"
)

var Current models.Config // Initialize the config struct

var (
	Debug      *bool                                                             // Global debug flag
	ConfigPath *string                                                           // Path to configuration file
	Password   *string                                                           // Password for authentication
	ServerPort *string                                                           // Port for server
	Secret     = make([]byte, 32)                                                // JWT secret key
	Submit     func(string, string, []string) ([]models.ResponseProtocol, error) // Function to submit data
	Cache      = true                                                            // Cache static file like css/js/image (If cache is enable more ram is used [default:true])
)

const (
	DefaultLimit   int = 100 // Default maximum number of flags to retrieve in the view
	DefaultOffeset int = 0   // Default offset for pagination
)
