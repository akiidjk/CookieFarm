package models

type Flag struct {
	SubmitTime   uint64 `json:"submit_time"`
	ResponseTime uint64 `json:"response_time"`
	ServicePort  uint16 `json:"service_port"`
	TeamID       uint16 `json:"team_id"`
	FlagCode     string `json:"flag_code"`
	ServiceName  string `json:"service_name"`
	Status       string `json:"status"`
}

type ResponseProtocol struct {
	Flag   string `json:"flag"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

type Service struct {
	Port uint16 `json:"port"`
	Name string `json:"name"`
}

type ConfigServer struct {
	SubmitFlagCheckerTime uint64 `json:"submit_flag_checker_time"`
	HostFlagchecker       string `json:"host_flagchecker"`
	TeamToken             string `json:"team_token"`
	MaxFlagBatchSize      uint16 `json:"max_flag_batch_size"`
	Protocol              string `json:"protocol"`
}

type ConfigClient struct {
	BaseUrlServer        string    `json:"base_url_server"`
	SubmitFlagServerTime uint64    `json:"submit_flag_server_time"`
	Services             []Service `json:"services"`
	RangeIPTeams         uint8     `json:"range_ip_teams"`
	FormatIPTeams        string    `json:"format_ip_teams"`
	MyTeamIP             string    `json:"my_team_ip"`
	RegexFlag            string    `json:"regex_flag"`
}

type Config struct {
	Configured   bool         `json:"configured"`
	ConfigServer ConfigServer `json:"server"`
	ConfigClient ConfigClient `json:"client"`
}

type Args struct {
	ExploitName   *string `json:"exploit_name"`
	Password      *string `json:"password"`
	BaseURLServer *string `json:"base_url_server"`
	TickTime      *int    `json:"tick_time"`
	Debug         *bool   `json:"debug"`
	Detach        *bool   `json:"detach"`
	ThreadCount   *int    `json:"thread_count"`
}

type TokenResponse struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

type ParsedFlagOutput struct {
	TeamID      uint16 `json:"team_id"`
	ServicePort uint16 `json:"service_port"`
	FlagCode    string `json:"flag_code"`
}
