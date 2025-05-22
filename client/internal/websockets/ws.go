// Package websockets used for communicating with the server via WebSocket protocol
package websockets

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/models"
	"github.com/gorilla/websocket"
)

var (
	OnNewConfig func()
)

const (
	FlagEvent   = "flag"
	ConfigEvent = "config"

	maxRetries = 3

	// Connection timeouts
	dialTimeout  = 10 * time.Second // Timeout for WebSocket connection
	writeTimeout = 10 * time.Second // Timeout for writing message

	StateClosed   CircuitState = iota // Allowed connections
	StateHalfOpen                     // Allowed connections with limited retries
	StateOpen                         // Blocked connections

	pongWait = 60 * time.Second // pongWait is the time to wait for a pong response
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// NewMessageEvent represents a new message event
type NewMessageEvent struct {
	Sent time.Time `json:"sent"`
}

// bad handshake (401)
// connection refused (503)
func GetConnection() (*websocket.Conn, error) {
	var conn *websocket.Conn
	var err error

	monitor := GetMonitor()
	monitor.SetStatus(StatusConnecting)

	if !circuitBreaker.IsAllowed() {
		logger.Log.Warn().Msg("Circuit breaker is open, connection attempt blocked")
		monitor.SetStatus(StatusFailed)
		return nil, ErrCircuitOpen
	}

	maxAttempts := maxRetries
	if circuitBreaker.state == StateHalfOpen {
		maxAttempts = halfOpenMaxRetry
	}

	dialer := &websocket.Dialer{
		HandshakeTimeout: dialTimeout,
	}

	for attempts := 0; attempts < maxAttempts; attempts++ {
		conn, _, err = dialer.Dial("ws://"+config.HostServer+"/ws", http.Header{
			"Cookie": []string{"token=" + config.Token},
		})

		if err == nil {
			circuitBreaker.RecordSuccess()

			monitor.SetConnection(conn)

			conn.SetPongHandler(func(appData string) error {
				monitor.mutex.Lock()
				monitor.stats.LastPongTime = time.Now()
				monitor.mutex.Unlock()
				return nil
			})

			conn.SetReadDeadline(time.Now().Add(pongWait))

			return conn, nil
		}

		if websocket.ErrBadHandshake == err {
			logger.Log.Error().Err(err).Msg("Bad handshake, retrying login...")
			config.Token, err = api.Login(config.Args.Password)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to refresh token")
				circuitBreaker.RecordFailure()
				monitor.SetStatus(StatusFailed)
				return nil, err
			}
			continue
		}

		logger.Log.Warn().Err(err).Int("attempt", attempts+1).Int("maxRetries", maxAttempts).Msg("Error connecting to WebSocket, retrying...")
		monitor.SetStatus(StatusReconnecting)
		time.Sleep(time.Second * time.Duration(1<<attempts))
	}

	circuitBreaker.RecordFailure()
	monitor.SetStatus(StatusFailed)
	monitor.RecordDisconnect(err)
	logger.Log.Error().Err(err).Msg("Failed to connect to WebSocket after multiple attempts")
	return nil, err
}

func WSReader(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error reading message from WebSocket")
			circuitBreaker.RecordFailure()
			break
		}

		if err := WSHandleMessage(message); err != nil {
			logger.Log.Error().Err(err).Msg("Error handling message")
			break
		}
	}
}

func WSHandleMessage(message []byte) error {
	var event Event
	if err := json.Unmarshal(message, &event); err != nil {
		return err
	}

	logger.Log.Debug().Str("type", event.Type).Str("Payload", string(event.Payload)).Msg("Received event")

	switch event.Type {
	case ConfigEvent:
		return ConfigHandler(event.Payload)
	default:
		//
	}

	return nil
}

func ConfigHandler(payload json.RawMessage) error {
	var configReceived models.Config
	if err := json.Unmarshal(payload, &configReceived); err != nil {
		return err
	}

	config.Current = configReceived

	if OnNewConfig != nil {
		go OnNewConfig()
	}

	return nil
}
