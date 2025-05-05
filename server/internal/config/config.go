// Config package for configuration management
package config

import (
	"github.com/ByteTheCookies/cookieserver/internal/models"
)

var Current models.Config // Initialize the config struct

var Debug *bool                                                              // Global debug flag
var ConfigPath *string                                                       // Path to configuration file
var Password *string                                                         // Password for authentication
var ServerPort *string                                                       // Port for server
var Secret = make([]byte, 32)                                                // JWT secret key
var Submit func(string, string, []string) ([]models.ResponseProtocol, error) // Function to submit data

const DEFAULT_LIMIT int = 100 // Default maximum number of flags to retrieve in the view
const DEFAULT_OFFSET int = 0  // Default offset for pagination
