package models

type Flag struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	FlagCode     string `json:"flag_code"`
	ServiceName  string `json:"service_name"`
	ServicePort  uint16 `json:"service_port"`
	TeamID       uint16 `json:"team_id"`
	SubmitTime   uint64 `json:"submit_time"`
	ResponseTime uint64 `json:"response_time"`
}

type ResponseProtocol struct {
	Msg    string `json:"msg"`
	Flag   string `json:"flag"`
	Status string `json:"status"`
}

type Service struct {
	Name string `json:"name"`
	Port uint16 `json:"port"`
}

type ConfigServer struct {
	CycleTime        uint64 `json:"cycle_time"`       // intervallo per invio flag al flagchecker
	HostFlagchecker  string `json:"host_flagchecker"` // es: localhost:3000
	TeamToken        string `json:"team_token"`
	MaxFlagBatchSize uint16 `json:"max_flag_batch_size"`
	Protocol         string `json:"protocol"` // Name of SO file protocol without extension
}

type ConfigClient struct {
	BaseUrlServer string    `json:"base_url_server"` // es: http://localhost:8080
	CycleTime     uint64    `json:"cycle_time"`      // intervallo per invio flag al server
	Services      []Service `json:"services"`
	RangeIpTeams  string    `json:"range_ip_teams"`  // min-max
	FormatIpTeams string    `json:"format_ip_teams"` // 10.0.0.{}
	MyTeamIp      string    `json:"my_team_ip"`      // Your IP in the A/D
}

type Config struct {
	ConfigServer ConfigServer `json:"server"`
	ConfigClient ConfigClient `json:"client"`
}

type TokenResponse struct {
	Exp   int64  `json:"exp"`
	Token string `json:"token"`
}
