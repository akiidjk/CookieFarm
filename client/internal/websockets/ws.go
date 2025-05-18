// Package websockets used for communicating with the server via WebSocket protocol
package websockets

import (
	"net/http"
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/gorilla/websocket"
)

const (
	FlagEvent  = "flag"
	maxRetries = 5
)

// bad handshake (401)
// connection refused (503)
func GetConnection() (*websocket.Conn, error) {
	var conn *websocket.Conn
	var err error

	for attempts := 0; attempts < maxRetries; attempts++ {
		conn, _, err = websocket.DefaultDialer.Dial("ws://"+*config.HostServer+"/ws", http.Header{
			"Cookie": []string{"token=" + config.Token},
		})

		if err == nil {
			return conn, nil
		}

		if websocket.ErrBadHandshake == err {
			logger.Log.Error().Err(err).Msg("Bad handshake, retrying login...")
			config.Token, err = api.Login(*config.Args.Password)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to refresh token")
				return nil, err
			}
			continue
		}

		logger.Log.Warn().Err(err).Int("attempt", attempts+1).Int("maxRetries", maxRetries).Msg("Error connecting to WebSocket, retrying...")
		time.Sleep(time.Second * time.Duration(attempts+1)) // Exponential backoff
	}

	logger.Log.Error().Err(err).Msg("Failed to connect to WebSocket after multiple attempts")
	return nil, err
}
