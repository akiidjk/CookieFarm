package models

type Flag struct {
	SubmitTime   uint64 `json:"submit_time"`
	ResponseTime uint64 `json:"response_time"`

	ID          string `json:"id"`
	Status      string `json:"status"`
	FlagCode    string `json:"flag_code"`
	ServiceName string `json:"service_name"`

	ServicePort uint16 `json:"service_port"`
	TeamID      uint16 `json:"team_id"`
}

type ResponseProtocol struct {
	Msg    string `json:"msg"`
	Flag   string `json:"flag"`
	Status string `json:"status"`
}

type Service struct {
	Name string
	Port uint16
}

type ConfigServer struct {
	HostFlagchecker string `json:"host_flagchecker"` // example: localhost:8080
	TeamToken       string `json:"team_token"`
	CycleTime       uint64 `json:"cycle_time"` // Time interval for send flags to flagchecker
}

type ConfigClient struct {
	BaseUrlServer string    `json:"base_url_server"` // Url example: http://localhost:8080
	CycleTime     uint64    `json:"cycle_time"`      // Time interval for send flags to server
	Services      []Service `json:"services"`
}

type Config struct {
	Server ConfigServer `json:"server"`
	Client ConfigClient `json:"client"`
}

type TokenResponse struct {
	Exp   int64  `json:"exp"`
	Token string `json:"token"`
}
