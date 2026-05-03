package api

import (
	serverconfig "server/config"
	"server/database"
)

// ResponseFlags represents the response for the flags api
type ResponseFlags struct {
	Flags  []database.Flag `json:"flags"`
	Next   string          `json:"next"`
	Nflags int64           `json:"n_flags"`
}

type ResponseChartStats struct {
	TickSeries              []FlagTickPoint              `json:"tick_series"`
	ExploitShare            []FlagExploitShare           `json:"exploit_share"`
	ExploitTickSeries       []FlagExploitTickSeries      `json:"exploit_tick_series"`
	ExploitStatusPercentage []FlagExploitStatusBreakdown `json:"exploit_status_percentage"`
	TotalFlags              int64                        `json:"total_flags"`
}

type FlagTickPoint struct {
	Timestamp int64 `json:"timestamp"`
	Total     int64 `json:"total"`
	Queued    int64 `json:"queued"`
	Accepted  int64 `json:"accepted"`
	Denied    int64 `json:"denied"`
	Error     int64 `json:"error"`
	Invalid   int64 `json:"invalid"`
}

type FlagExploitShare struct {
	Name       string  `json:"name"`
	Value      int64   `json:"value"`
	Percentage float64 `json:"percentage"`
}

type FlagExploitTickSeries struct {
	Name  string                 `json:"name"`
	Total int64                  `json:"total"`
	Data  []FlagExploitTickPoint `json:"data"`
}

type FlagExploitTickPoint struct {
	Timestamp int64 `json:"timestamp"`
	Value     int64 `json:"value"`
}

type FlagExploitStatusBreakdown struct {
	Name          string  `json:"name"`
	Total         int64   `json:"total"`
	Queued        float64 `json:"queued"`
	Accepted      float64 `json:"accepted"`
	Denied        float64 `json:"denied"`
	Error         float64 `json:"error"`
	Invalid       float64 `json:"invalid"`
	QueuedCount   int64   `json:"queued_count"`
	AcceptedCount int64   `json:"accepted_count"`
	DeniedCount   int64   `json:"denied_count"`
	ErrorCount    int64   `json:"error_count"`
	InvalidCount  int64   `json:"invalid_count"`
}

// SigninRequest from the client to the server
type SigninRequest struct {
	Username string `json:"username,omitzero"` // Username for authentication
	Password string `json:"password"`
}

type AuthVerifyResponse struct {
	Username string `json:"username"`
}

// Pagination structure for manage data in the view
type Pagination struct {
	PageList []int // List of page numbers to display in the pagination
	Pages    int   // Total number of pages
	Limit    int   // Maximum number of items per page
	Current  int   // Current page number (offset / limit)
	HasPrev  bool  // Indicates if there is a previous page
	HasNext  bool  // Indicates if there is a next page
}

// ViewParamsDashboard represents the parameters for the dashboard view
type ViewParamsDashboard struct {
	Limit int `json:"limit"` // Maximum number of items per page
}

// ViewParamsPagination represents the parameters for the pagination view
type ViewParamsPagination struct {
	Pagination Pagination // Pagination parameters
}

// ViewParamsFlags represents the parameters for the flags view
type ViewParamsFlags struct {
	Flags []database.Flag `json:"flags"` // List of flags to display
}

// ResponseSuccess represents the response for the success api
type ResponseSuccess struct {
	Message string `json:"message"` // Message for the success response
}

// ResponseError represents the response for the error api
type ResponseError struct {
	Error   string `json:"error"`   // Error message for the error response
	Details string `json:"details"` // Details for the error response
}

// UpdateConfigRequest wraps full configuration payload.
type UpdateConfigRequest struct {
	Config serverconfig.FullConfig `json:"config"`
}

// ResponseSharedConfig represents the full configuration returned by API.
type ResponseSharedConfig = serverconfig.FullConfig
