package protocols

// ResponseProtocol represents a response from the flag checker service.
type ResponseProtocol struct {
	Status int64  `json:"status"` // Status of the response (e.g., "0", "1", "2" ,"3") see enum in pkg/models/models.go
	Flag   string `json:"flag"`   // Flag string received from the flag checker service
	Msg    string `json:"msg"`    // Message from the flag checker service
}
