package core

import (
	"context"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/sqlite"
	"github.com/ByteTheCookies/cookieserver/protocols"
)

var shutdownCancel context.CancelFunc

// ----------- END FLAG GROUPS ------------

// StartFlagProcessingLoop starts the flag processing loop.
func StartFlagProcessingLoop(ctx context.Context) {
	interval := time.Duration(config.Current.ConfigServer.SubmitFlagCheckerTime) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Log.Info().Msg("Starting flag processing loop...")

	var err error
	config.Submit, err = protocols.LoadProtocol(config.Current.ConfigServer.Protocol)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to load protocol")
		return
	}

	// Main loop for flag processing.
	for {
		logger.Log.Debug().Msg("Waiting for flags to process...")
		select {
		case <-ctx.Done():
			logger.Log.Info().Msg("Flag processing loop terminated")
			return
		case <-ticker.C:
			flags, err := sqlite.GetUnsubmittedFlagCodeList(config.Current.ConfigServer.MaxFlagBatchSize)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to get unsubmitted flags")
				continue
			}
			if len(flags) == 0 {
				logger.Log.Debug().Msg("No flags to submit")
				continue
			}

			logger.Log.Info().Int("count", len(flags)).Msg("Submitting flags to checker")

			responses, err := config.Submit(
				config.Current.ConfigServer.HostFlagchecker,
				config.Current.ConfigServer.TeamToken,
				flags,
			)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Error submitting flags to checker")
				continue
			}

			UpdateFlags(responses)
		}
	}
}

// UpdateFlags updates the status of flags in the database.
func UpdateFlags(flags []protocols.ResponseProtocol) {
	valid := flags[:0]

	accepted, denied, errored := 0, 0, 0
	for _, f := range flags {
		switch f.Status {
		case sqlite.StatusAccepted:
			accepted++
			valid = append(valid, f)

		case sqlite.StatusDenied:
			denied++
			valid = append(valid, f)

		case sqlite.StatusError:
			errored++
			valid = append(valid, f)

		default:
			continue
		}
	}

	if err := sqlite.UpdateFlagsStatus(valid); err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Failed to update flags")
	}

	total := accepted + denied + errored
	logger.Log.Info().
		Int("accepted", accepted).
		Int("denied", denied).
		Int("errored", errored).
		Int("total", total).
		Msg("Flags update summary")
}
