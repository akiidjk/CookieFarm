package submitter

import (
	"time"

	"github.com/ByteTheCookies/cookiefarm-client/internal/api"
	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
)

func Start(flagsChan <-chan models.Flag) {
	submitInterval := time.Duration(config.Current.ConfigClient.SubmitFlagServerTime) * time.Second
	logger.Log.Debug().Dur("interval", submitInterval).Msg("Starting submitter loop")

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
			}
		}
	}
}
