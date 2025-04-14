package models

type Flag struct {
	ID           int64  `json:"id"`
	FlagCode     string `json:"flag_code"`
	ServiceName  string `json:"service_name"`
	SubmitTime   int64  `json:"submit_time"`
	ResponseTime int64  `json:"response_time"`
	Status       string `json:"status"`
	TeamID       int64  `json:"team_id"`
}

type Config struct {
}

type FlagResponse struct {
	Flags []Flag `json:"flags"`
}
