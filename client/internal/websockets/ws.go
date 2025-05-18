// Package websockets used for communicating with the server via WebSocket protocol
package websockets

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/gorilla/websocket"
)

const (
	FlagEvent  = "flag"
	maxRetries = 3
	// Circuit breaker constants
	failureThreshold = 2                // Numero di errori consecutivi prima di aprire il circuito
	resetTimeout     = 30 * time.Second // Tempo di attesa prima di tentare una nuova connessione
	halfOpenMaxRetry = 1                // Tentativi durante lo stato half-open
	// Connection timeouts
	dialTimeout  = 10 * time.Second // Timeout per la connessione WebSocket
	writeTimeout = 10 * time.Second // Timeout per la scrittura dei messaggi
)

// CircuitState state of the circuit breaker
type CircuitState int

const (
	StateClosed   CircuitState = iota // Allowed connections
	StateHalfOpen                     // Allowed connections with limited retries
	StateOpen                         // Blocked connections
)

// CircuitBreaker is a struct that implements the circuit breaker pattern
type CircuitBreaker struct {
	state           CircuitState
	failureCount    int
	lastFailureTime time.Time
	mutex           sync.Mutex
}

var (
	circuitBreaker = &CircuitBreaker{ // circuitBreaker is the instance of the circuit breaker
		state:        StateClosed,
		failureCount: 0,
	}

	// Errors
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

const (
	pongWait = 60 * time.Second // pongWait is the time to wait for a pong response

)

// CircuitBreaker methods

// RecordSuccess registers a successful connection and resets the failure count
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failureCount = 0
	cb.state = StateClosed
	logger.Log.Debug().Msg("Circuit breaker: Connection successful, circuit closed")
}

// RecordFailure registers a failed connection attempt and updates the state
func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.state == StateHalfOpen || cb.failureCount >= failureThreshold {
		cb.state = StateOpen
		logger.Log.Warn().Int("failures", cb.failureCount).Msg("Circuit breaker opened")
	}
}

// IsAllowed checks if a connection attempt is allowed based on the circuit breaker state
func (cb *CircuitBreaker) IsAllowed() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == StateClosed {
		return true
	}

	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > resetTimeout {
			cb.state = StateHalfOpen
			logger.Log.Info().Msg("Circuit breaker: Switched to half-open state")
			return true
		}
		return false
	}

	return true
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
		conn, _, err = dialer.Dial("ws://"+*config.HostServer+"/ws", http.Header{
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
			config.Token, err = api.Login(*config.Args.Password)
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
