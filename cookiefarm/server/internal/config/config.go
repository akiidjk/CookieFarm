// Package config for configuration management
package config

import (
	"protocols"
	"sharedconfig"
)

type Config struct {
	flagTTL uint64
}

var (
	SharedConfig sharedconfig.Shared // Initialize the config struct
	Config       Config
)

var (
	Debug         bool                                                                 // Global debug flag
	UseConfigFile bool                                                                 // Use configuration file instead of web config
	Password      string                                                               // Password for authentication
	ServerPort    string                                                               // Port for server
	Secret        = make([]byte, 32)                                                   // JWT secret key
	Submit        func(string, string, []string) ([]protocols.ResponseProtocol, error) // Function to submit data
	Cache         = true                                                               // Cache static file like css/js/image (If cache is enable more ram is used [default:true])
)

const (
	ConfigPath    string = "config.yml" // Path to configuration file
	DefaultLimit  int    = 100          // Default maximum number of flags to retrieve in the view
	DefaultOffset int    = 0            // Default offset for pagination
)
