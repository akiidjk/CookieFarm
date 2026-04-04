// Package config for configuration management
package config

import (
	"protocols"
	"sharedconfig"
	"sync"
)

type Config struct {
<<<<<<< Updated upstream:cookiefarm/server/internal/config/config.go
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
=======
	URLFlagChecker        string `json:"url_flag_checker" yaml:"url_flag_checker"`                 // Address of the flagchecker server
	TeamToken             string `json:"team_token" yaml:"team_token"`                             // Authentication token for team identity
	Protocol              string `json:"protocol" yaml:"protocol"`                                 // Protocol used to communicate with the flagchecker server
	StartTime             string `json:"start_time" yaml:"start_time"`                             // CTF competition start time (HH:MM:SS format)
	EndTime               string `json:"end_time" yaml:"end_time"`                                 // CTF competition end time (HH:MM:SS format)
	MaxFlagBatchSize      uint   `json:"max_flag_batch_size" yaml:"max_flag_batch_size"`           // Max number of flags to send in a single batch
	TickTime              int    `json:"tick_time" yaml:"tick_time"`                               // Duration of one game tick in seconds
	SubmitFlagCheckerTime uint64 `json:"submit_flag_checker_time" yaml:"submit_flag_checker_time"` // Time interval (s) to check and submit flags
	FlagTTL               uint64 `json:"flag_ttl" yaml:"flag_ttl"`                                 // Time-to-live for flags in ticks
}

var cfg Config // Global configurations
>>>>>>> Stashed changes:cookiefarm/server/config/config.go

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
