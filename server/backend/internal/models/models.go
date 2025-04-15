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

type Config struct {
}
