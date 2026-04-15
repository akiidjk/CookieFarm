package models

import "server/database"

const (
	StatusUnsubmitted = iota // Status for unsubmitted flags
	StatusAccepted           // Status for accepted flags
	StatusDenied             // Status for denied flags
	StatusError              // Status for error flags
	StatusNotValid
)

// SubmitFlagsRequest the struct for the requests from the client to server
type SubmitFlagsRequest struct {
	Flags []database.Flag `json:"flags"` // Flags to submit
}

// SubmitFlagRequest the struct for the requests from the client to server
type SubmitFlagRequest struct {
	Flag database.Flag `json:"flag"` // Flag to submit
}
