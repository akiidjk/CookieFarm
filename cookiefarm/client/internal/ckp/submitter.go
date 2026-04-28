package ckp

import (
	"logger"
	"server/database"
)

const ADDR = "127.0.0.1:7777"

func Start(flagsChan <-chan database.Flag) {
	logger.Log.Debug().Msg("Starting submission loop to the cookiefarm server with websockets (websockets) ...")

	conn, err := NewClient(ADDR)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error connecting to WebSocket")
	}
	defer conn.Close()

	logger.Log.Info().Msg("Successfully connected to CKP server")
	for flag := range flagsChan {
		flagBytes := buildPayload(flag)
		conn.SendWithRetry(ADDR, flagBytes, 3)
	}

	logger.Log.Info().Msg("Submission loop finished")
}
