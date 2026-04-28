package ckp

import (
	"logger"
	"server/database"
)

const ADDR = "127.0.0.1:7777"

func Start(flagsChan <-chan database.Flag) {
	logger.Log.Debug().Msg("Starting submission loop to the cookiefarm server with ckp ...")

	conn, err := NewClient(ADDR)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error connecting to WebSocket")
	}
	defer conn.Close()
	go conn.ReadPump()

	for flag := range flagsChan {
		flagBytes := buildPayload(flag)
		err := conn.SendWithRetry(ADDR, flagBytes, 3)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error sending flag to CKP server")
		}
	}

	logger.Log.Info().Msg("Submission loop finished")
}
