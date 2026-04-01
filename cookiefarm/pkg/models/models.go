package models

import "server/database"

const (
	StatusUnsubmitted = iota // Status for unsubmitted flags
	StatusAccepted           // Status for accepted flags
	StatusDenied             // Status for denied flags
	StatusError              // Status for error flags
	StatusNotValid
)

type Service struct {
	Name string `json:"name" yaml:"name"` // Name identifier of the service
	Port uint16 `json:"port" yaml:"port"` // Port where the service is exposed
}

// SubmitFlagsRequest the struct for the requests from the client to server
type SubmitFlagsRequest struct {
	Flags []database.Flag `json:"flags"` // Flags to submit
}

// SubmitFlagRequest the struct for the requests from the client to server
type SubmitFlagRequest struct {
	Flag database.Flag `json:"flag"` // Flag to submit
}
