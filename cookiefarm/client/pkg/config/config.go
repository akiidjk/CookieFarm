package config

import (
	"sync"
	"system"
)

var DefaultPath, _ = system.ExpandTilde("~/.config/cookiefarm")

type InfraMetaData struct {
	FlagRegex     string `json:"regex_flag" yaml:"regex_flag"`           // Regex used to identify flags in output
	FormatIPTeams string `json:"format_ip_teams" yaml:"format_ip_teams"` // Format string for generating team IPs
	MyTeamID      int    `json:"my_team_id" yaml:"my_team_id"`           // ID of the current team
	URLFlagIds    string `json:"url_flag_ids" yaml:"url_flag_ids"`       // URLFlagIds is the where the flagsId server is running
	NOPTeam       int    `json:"nop_team" yaml:"nop_team"`               // The id of the nop team in the ctf
	RangeIPTeams  uint8  `json:"range_ip_teams" yaml:"range_ip_teams"`   // Number of teams / IP range
}

type Config struct {
	Host     string `json:"host" yaml:"host"`         // Host address of the server
	Username string `json:"username" yaml:"username"` // Username of the client
	Port     uint16 `json:"port" yaml:"port"`         // Port of the server
	HTTPS    bool   `json:"protocol" yaml:"protocol"` // Protocol used to connect to the server (e.g., http, https)

	services map[string]uint16
}

type ConfigManager struct {
	mu    sync.RWMutex
	cfg   Config
	token string
}
