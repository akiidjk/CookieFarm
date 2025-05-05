// Package submitter provides functions to manage the CookieFarm client submission.
package submitter

import (
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/api"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/models"
)

// Start initializes the submission loop to the cookiefarm server.
func Start(flagsChan <-chan models.Flag) {
	submitInterval := time.Duration(config.Current.ConfigClient.SubmitFlagServerTime) * time.Second
	logger.Log.Info().Uint64("interval seconds", config.Current.ConfigClient.SubmitFlagServerTime).Msg("Starting submission loop...")

	ticker := time.NewTicker(submitInterval)
	defer ticker.Stop()

	var flagsBatch []models.Flag

	for {
		select {
		case flag := <-flagsChan:
			flagsBatch = append(flagsBatch, flag)
		case <-ticker.C:
			if len(flagsBatch) > 0 {
				err := api.SendFlag(flagsBatch...)
				if err != nil {
					logger.Log.Error().Err(err).Msg("Error sending flags batch")
				}
				logger.Log.Info().Int("flags_sent", len(flagsBatch)).Msg("Submitted flags batch")
				flagsBatch = nil
			} else {
				logger.Log.Info().Msg("No flags to send")
			}
		}
	}
}
