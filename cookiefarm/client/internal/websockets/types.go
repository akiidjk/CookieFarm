package websockets

import (
	"encoding/json"
	"server/database"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type CircuitState int

type CircuitBreaker struct {
	state           CircuitState
	failureCount    int
	lastFailureTime time.Time
	mutex           sync.Mutex
}

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type NewMessageEvent struct {
	Sent time.Time `json:"sent"`
}

type EventWS struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

type EventWSFlag struct {
	Type    string        `json:"type"`
	Payload database.Flag `json:"payload"`
}

// ConnectionStatus represents the current status of the WebSocket connection
type (
	ConnectionStatus   int
	ConnectionTracking struct {
		ConnectionAttempts int
		SuccessfulConnects int
		FailedConnects     int
		LastConnectTime    time.Time
		LastDisconnectTime time.Time
	}
	ConnectionMessages struct {
		MessagesSent     int
		MessagesReceived int
		LastSendTime     time.Time
		LastReceiveTime  time.Time
	}
	LatencyTracking struct {
		LastPingTime     time.Time
		LastPongTime     time.Time
		CurrentLatency   time.Duration
		AverageLatency   time.Duration
		totalLatencySum  time.Duration
		latencyDataCount int
	}
)

const (
	StatusDisconnected ConnectionStatus = iota
	StatusConnecting
	StatusConnected
	StatusReconnecting
	StatusFailed
)

type ConnectionStats struct {
	ConnectionTracking ConnectionTracking
	MessageTracking    ConnectionMessages
	LastError          error
	ConsecutiveErrs    int
	CurrentStatus      ConnectionStatus
	LatencyTracking    LatencyTracking
}

type ConnectionMonitor struct {
	stats        ConnectionStats
	mutex        sync.RWMutex
	conn         *websocket.Conn
	isMonitoring bool
	stopChan     chan struct{}
}
