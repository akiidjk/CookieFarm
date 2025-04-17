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
	HostFlagchecker string `json:"host_flagchecker"` // es: localhost:8080
	TeamToken       string `json:"team_token"`
	CycleTime       uint64 `json:"cycle_time"` // intervallo per invio flag al flagchecker
}

type ConfigClient struct {
	BaseUrlServer string    `json:"base_url_server"` // es: http://localhost:8080
	CycleTime     uint64    `json:"cycle_time"`      // intervallo per invio flag al server
	Services      []Service `json:"services"`
}

type Config struct {
	ConfigServer ConfigServer `json:"server"`
	ConfigClient ConfigClient `json:"client"`
}

type TokenResponse struct {
	Exp   int64  `json:"exp"`
	Token string `json:"token"`
}
