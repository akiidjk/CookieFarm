package sqlite

const (
	StatusUnsubmitted = "UNSUBMITTED" // Status for unsubmitted flags
	StatusAccepted    = "ACCEPTED"    // Status for accepted flags
	StatusDenied      = "DENIED"      // Status for denied flags
	StatusError       = "ERROR"       // Status for error flags
)

// Flag represents a single flag captured during a CTF round.
// It includes metadata about the submission and the service context.
type Flag struct {
	SubmitTime   uint64 `json:"submit_time"`   // UNIX timestamp when the flag was submitted
	ResponseTime uint64 `json:"response_time"` // UNIX timestamp when a response was received
	FlagCode     string `json:"flag_code"`     // Actual flag string
	ServiceName  string `json:"service_name"`  // Human-readable name of the service
	Status       string `json:"status"`        // Status of the submission (e.g., "unsubmitted", "accepted", "denied")
	Msg          string `json:"msg"`           // Message from the flag checker service
	PortService  uint16 `json:"port_service"`  // Port of the vulnerable service
	TeamID       uint16 `json:"team_id"`       // ID of the team the flag was captured from
}
