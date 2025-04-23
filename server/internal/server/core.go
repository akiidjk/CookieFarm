package server

import (
	"context"
	"time"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/protocols"
)

func (s *FiberServer) StartFlagProcessingLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(config.Current.ConfigServer.SubmitFlagCheckerTime) * time.Second)
	defer ticker.Stop()

	logger.Log.Info().Msg("Starting flag processing loop")

	var err error
	config.Submit, err = protocols.LoadProtocol(config.Current.ConfigServer.Protocol)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to load protocol")
		return
	}

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info().Msg("Flag processing loop terminated")
			return

		case <-ticker.C:
			flags, err := s.db.GetUnsubmittedFlagCodeList(int(config.Current.ConfigServer.MaxFlagBatchSize))
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to get unsubmitted flags")
				continue
			}

			if len(flags) == 0 {
				logger.Log.Debug().Msg("No flags to submit at this time")
				continue
			}

			logger.Log.Debug().Int("count", len(flags)).Msg("Submitting flags to checker")

			responses, err := config.Submit(
				config.Current.ConfigServer.HostFlagchecker,
				config.Current.ConfigServer.TeamToken,
				flags,
			)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Error submitting flags to checker")
				continue
			}

			s.UpdateFlags(responses)
		}
	}
}

func (s *FiberServer) UpdateFlags(flags []models.ResponseProtocol) {
	var (
		flagsAccepted []string
		flagsDenied   []string
		flagsErrored  []string
	)

	for _, resp := range flags {
		switch resp.Status {
		case "ACCEPTED":
			flagsAccepted = append(flagsAccepted, resp.Flag)
			logger.Log.Debug().Str("flag", resp.Flag).Msg("Flag accepted")
		case "DENIED":
			flagsDenied = append(flagsDenied, resp.Flag)
			logger.Log.Debug().Str("flag", resp.Flag).Msg("Flag denied")
		default:
			flagsErrored = append(flagsErrored, resp.Flag)
			logger.Log.Debug().Str("flag", resp.Flag).Msg("Flag error")
		}
	}

	if len(flagsAccepted) > 0 {
		if err := s.db.UpdateFlagsStatus(flagsAccepted, "ACCEPTED"); err != nil {
			logger.Log.Error().Err(err).Int("count", len(flagsAccepted)).Msg("Failed to update accepted flags")
		}
	}
	if len(flagsDenied) > 0 {
		if err := s.db.UpdateFlagsStatus(flagsDenied, "DENIED"); err != nil {
			logger.Log.Error().Err(err).Int("count", len(flagsDenied)).Msg("Failed to update denied flags")
		}
	}
	if len(flagsErrored) > 0 {
		if err := s.db.UpdateFlagsStatus(flagsErrored, "ERROR"); err != nil {
			logger.Log.Error().Err(err).Int("count", len(flagsErrored)).Msg("Failed to update errored flags")
		}
	}

	logger.Log.Info().
		Int("accepted", len(flagsAccepted)).
		Int("denied", len(flagsDenied)).
		Int("error", len(flagsErrored)).
		Int("total", len(flagsAccepted)+len(flagsDenied)+len(flagsErrored)).
		Msg("Flags update summary")
}
