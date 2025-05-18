// Package submitter provides functions to manage the CookieFarm client submission.
package submitter

import (
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/models"
	"github.com/ByteTheCookies/cookieclient/internal/websockets"
	json "github.com/bytedance/sonic"
	gorilla "github.com/gorilla/websocket"
)

// Start initializes the submission loop to the cookiefarm server.
func Start(flagsChan <-chan models.Flag) error {
	logger.Log.Info().Msg("Starting submission loop to the cookiefarm server...")
	conn, err := websockets.GetConnection()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error connecting to WebSocket")
	}
	defer conn.Close()

	for {
		select {
		case flag := <-flagsChan:
			flagObj := models.EventWSFlag{
				Type:    websockets.FlagEvent,
				Payload: flag,
			}
			marshalFlag, err := json.Marshal(flagObj)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Error marshalling flag")
				continue
			}
			if err := conn.WriteMessage(gorilla.TextMessage, marshalFlag); err != nil {
				logger.Log.Error().Err(err).Msg("Error sending flag, attempting reconnection")
				newConn, reconnectErr := websockets.GetConnection()
				if reconnectErr != nil {
					logger.Log.Fatal().Err(reconnectErr).Msg("Failed to reconnect to WebSocket")
				} else {
					conn = newConn
					logger.Log.Info().Msg("Successfully reconnected to WebSocket")
				}
				continue
			}
		}
	}
}

// TODO: Sistema di accumulo in caso di non connesione fino a 2 minuti
