// Config package for configuration management
package config

import (
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/ByteTheCookies/cookieserver/internal/utils"
)

var Current models.Config = models.Config{ // Initialize the config struct
	Configured:   false,
	ConfigClient: models.ConfigClient{},
	ConfigServer: models.ConfigServer{},
}

var Debug *bool                                                              // Global debug flag
var ConfigPath string                                                        // Path to configuration file
var Secret = make([]byte, 32)                                                // JWT secret key
var Password = utils.GetEnv("PASSWORD", "password")                          // Password for authentication
var ServerPort = utils.GetEnv("BACKEND_PORT", "8080")                        // Port for server
var Submit func(string, string, []string) ([]models.ResponseProtocol, error) // Function to submit data

const DEFAULT_LIMIT = 100 // Default maximum number of flags to retrieve in the view
const DEFAULT_OFFSET = 0  // Default offset for pagination
