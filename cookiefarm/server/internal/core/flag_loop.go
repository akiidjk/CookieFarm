package core

import (
	"context"
	"logger"
	"models"
	"protocols"
	"time"

	"server/config"
	"server/database"
)

// ----------- END FLAG GROUPS ------------

// StartFlagProcessingLoop starts the flag processing loop.
func (r *Runner) StartFlagProcessingLoop(ctx context.Context) {
	interval := time.Duration(r.config.GetSubmitFlagCheckerTime()) * time.Second
	if interval <= 0 {
		logger.Log.Warn().Msgf("Invalid SubmitFlagCheckerTime %d, defaulting to 60 seconds", r.config.GetSubmitFlagCheckerTime())
		interval = 60 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Log.Info().Msg("Starting flag processing loop...")

	var err error
	if config.Submit == nil {
		config.Submit, err = protocols.LoadProtocol(r.config.GetProtocol()) // Change from config
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to load protocol")
			return
		}
	}

	// Main loop for flag processing.
	for {
		logger.Log.Debug().Msg("Waiting for flags to process...")
		select {
		case <-ctx.Done():
			logger.Log.Info().Msg("Flag processing loop terminated")
			return
		case <-ticker.C:
			flags, err := r.store.Queries.GetUnsubmittedFlagCodes(ctx, int64(r.config.GetMaxFlagBatchSize())) // Cast not good, but we know the value is within int64 range
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
				r.config.GetURLFlagChecker(),
				r.config.GetTeamToken(),
				flags,
			)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Error submitting flags to checker")
				continue
			}
			r.UpdateFlags(responses)
		}
	}
}

func (r *Runner) UpdateFlags(flags []protocols.ResponseProtocol) {
	statusCounts := map[int64]int{
		models.StatusAccepted: 0,
		models.StatusDenied:   0,
		models.StatusError:    0,
	}

	valid := flags[:0] // Reuse the same slice to avoid extra allocations
	logger.Log.Debug().Msgf("Processing %d flag responses", len(flags))
	for _, f := range flags {
		if _, exists := statusCounts[f.Status]; exists {
			statusCounts[f.Status]++
			valid = append(valid, f)
		}
	}

	ctx := context.Background()

	err := r.store.WithTx(ctx, func(q *database.Queries) error {
		for _, f := range valid {
			if err := q.UpdateFlagStatusByCode(ctx, database.MapFromResponseProtocolToParamsToUpdate(f)); err != nil {
				logger.Log.Error().
					Err(err).
					Str("flag_code", f.Flag).
					Int64("status", f.Status).
					Msg("Failed to update flag status")
			}
		}
		return nil
	})
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to update flags in database")
		return
	}

	total := statusCounts[models.StatusAccepted] + statusCounts[models.StatusDenied] + statusCounts[models.StatusError]
	logger.Log.Info().
		Int("accepted", statusCounts[models.StatusAccepted]).
		Int("denied", statusCounts[models.StatusDenied]).
		Int("errored", statusCounts[models.StatusError]).
		Int("total", total).
		Msg("Flags update summary")
}
