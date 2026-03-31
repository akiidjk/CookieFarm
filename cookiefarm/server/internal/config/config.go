// Package config for configuration management
package config

import (
	"protocols"
	"sharedconfig"
	"sync"
)

type Config struct {
	URLFlagChecker        string `yaml:"url_flag_checker" json:"url_flag_checker"`
	TeamToken             string `yaml:"team_token" json:"team_token"`
	SubmitFlagCheckerTime uint   `yaml:"submit_flag_checker_time" json:"submit_flag_checker_time"`
	MaxFlagBatchSize      uint   `yaml:"max_flag_batch_size" json:"max_flag_batch_size"`
	Protocol              string `yaml:"protocol" json:"protocol"`
	TickTime              uint   `yaml:"tick_time" json:"tick_time"`
	FlagTTL               uint64 `yaml:"flag_ttl" json:"flag_ttl"`
	StartTime             string `yaml:"start_time" json:"start_time"`
	EndTime               string `yaml:"end_time" json:"end_time"`
}

type FullConfig struct {
	Server     Config              `yaml:"server" json:"server"`
	Shared     sharedconfig.Shared `yaml:"shared" json:"shared"`
	Configured bool                `yaml:"configured" json:"configured"`
}

type ConfigManager struct {
	mu    sync.RWMutex
	cfg   FullConfig
	token string
}

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
