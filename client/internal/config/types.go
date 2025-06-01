package config

// Service represents a single vulnerable service as defined in the configuration.
type Service struct {
	Port uint16 `json:"port" yaml:"port"` // Port where the service is exposed
	Name string `json:"name" yaml:"name"` // Name identifier of the service
}

// ConfigServer holds configuration data required by the server to submit and validate flags.
type ConfigServer struct {
	SubmitFlagCheckerTime uint64 `json:"submit_flag_checker_time" yaml:"submit_flag_checker_time"` // Time interval (s) to check and submit flags
	MaxFlagBatchSize      uint   `json:"max_flag_batch_size" yaml:"max_flag_batch_size"`           // Max number of flags to send in a single batch
	HostFlagchecker       string `json:"host_flagchecker" yaml:"host_flagchecker"`                 // Address of the flagchecker server
	TeamToken             string `json:"team_token" yaml:"team_token"`                             // Authentication token for team identity
	Protocol              string `json:"protocol" yaml:"protocol"`                                 // Protocol used to communicate with the flagchecker server
}

// ConfigClient contains all client-side configuration options.
type ConfigClient struct {
	RangeIPTeams  uint8     `json:"range_ip_teams" yaml:"range_ip_teams"`   // Number of teams / IP range
	Services      []Service `json:"services" yaml:"services"`               // List of services to exploit
	FormatIPTeams string    `json:"format_ip_teams" yaml:"format_ip_teams"` // Format string for generating team IPs
	MyTeamIP      string    `json:"my_team_ip" yaml:"my_team_ip"`           // IP address of the current team
	RegexFlag     string    `json:"regex_flag" yaml:"regex_flag"`           // Regex used to identify flags in output
}

// Config aggregates both server and client configuration,
// and includes a flag indicating whether the configuration is initialized.
type Config struct {
	Configured   bool         `json:"configured" yaml:"configured"` // True if configuration has been loaded and validated
	ConfigServer ConfigServer `json:"server" yaml:"server"`         // Server-specific configuration
	ConfigClient ConfigClient `json:"client" yaml:"client"`         // Client-specific configuration
}

// ArgsAttack represents the command-line arguments or configuration
// values passed at runtime to control the exploit manager behavior.
type ArgsAttack struct {
	ServicePort uint16 `json:"port"`         // Service Port
	TickTime    int    `json:"tick_time"`    // Optional custom tick time
	ThreadCount int    `json:"thread_count"` // Optional number of concurrent threads (coroutine) to use
	Debug       bool   `json:"debug"`        // Enable debug mode if true
	Detach      bool   `json:"detach"`       // Run in background if true
	ExploitPath string `json:"exploit_path"` // Path to the exploit to run
}

type Exploit struct {
	Name string `json:"name"` // Name of the exploit
	PID  int    `json:"pid"`  // Process ID of the exploit
}

type ArgsConfig struct {
	Address  string    `json:"address"`  // Host address of the server
	Port     uint16    `json:"port"`     // Port of the server
	HTTPS    bool      `json:"protocol"` // Protocol used to connect to the server (e.g., http, https)
	Username string    `json:"username"` // Username of the client
	Exploits []Exploit `json:"exploits"` // List of exploits available in the client
}
