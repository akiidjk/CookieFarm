package models

type Flag struct {
	SubmitTime   uint64 `json:"submit_time"`
	ResponseTime uint64 `json:"response_time"`
	ServicePort  uint16 `json:"service_port"`
	TeamID       uint16 `json:"team_id"`
	ID           string `json:"id"`
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

func (s ConfigServer) IsEmpty() bool {
	return s.HostFlagchecker == "" && s.TeamToken == "" && s.SubmitFlagCheckerTime == 0
}

type ConfigClient struct {
	BaseUrlServer        string    `json:"base_url_server"`
	SubmitFlagServerTime uint64    `json:"submit_flag_server_time"`
	Services             []Service `json:"services"`
	TeamsIPRange         uint8     `json:"range_ip_teams"`
	TeamIPFormat         string    `json:"format_ip_teams"`
	MyTeamIP             string    `json:"my_team_ip"`
	RegexFlag            string    `json:"regex_flag"`
}

func (c ConfigClient) IsEmpty() bool {
	return c.BaseUrlServer == "" && c.SubmitFlagServerTime == 0 && len(c.Services) == 0
}

type Config struct {
	Configured   bool         `json:"configured"`
	ConfigServer ConfigServer `json:"server"`
	ConfigClient ConfigClient `json:"client"`
}

func (c Config) IsEmpty() bool {
	return c.ConfigServer.IsEmpty() && c.ConfigClient.IsEmpty()
}

type Session struct {
	Password string `json:"password"`
}

type SigninRequest struct {
	Password string `json:"password"`
}
