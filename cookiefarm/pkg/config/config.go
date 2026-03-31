package sharedconfig

import (
	"fmt"
	"models"
	"runtime/debug"
	"sync"
)

type Shared struct {
	Services      []models.Service `json:"services" yaml:"services"`               // List of services to exploit
	RegexFlag     string           `json:"regex_flag" yaml:"regex_flag"`           // Regex used to identify flags in output
	FormatIPTeams string           `json:"format_ip_teams" yaml:"format_ip_teams"` // Format string for generating team IPs
	MyTeamID      int              `json:"my_team_id" yaml:"my_team_id"`           // ID of the current team
	URLFlagIds    string           `json:"url_flag_ids" yaml:"url_flag_ids"`       // URLFlagIds is the where the flagsId server is running
	NOPTeam       int              `json:"nop_team" yaml:"nop_team"`               // The id of the nop team in the ctf
	RangeIPTeams  uint8            `json:"range_ip_teams" yaml:"range_ip_teams"`   // Number of teams / IP range
	Configured    bool             `json:"configured" yaml:"configured"`           // True if configuration has been loaded and validated
}

var (
	instance *Shared
	once     sync.Once
)

func GetInstance() *Shared {
	once.Do(func() {
		instance = &Shared{}
	})
	return instance
}

func (cfg *Shared) Set(new Shared) {
	*cfg = new
}

func GetVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Build info not available")
		return ""
	}
	return info.Main.Version
}
