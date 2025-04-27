package server

import (
	"context"
	"sync"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/ByteTheCookies/cookieserver/protocols"
)

// ----------- FLAG GROUPS ------------

type FLagGroups struct {
	accepted []string
	denied   []string
	errored  []string
	// capienza massima prevista
	cap int
}

func newFlagGroups(cap int) *FLagGroups {
	return &FLagGroups{
		accepted: make([]string, 0, cap),
		denied:   make([]string, 0, cap),
		errored:  make([]string, 0, cap),
		cap:      cap,
	}
}

func (g *FLagGroups) reset() {
	g.accepted = g.accepted[:0]
	g.denied = g.denied[:0]
	g.errored = g.errored[:0]
}

// ----------- END FLAG GROUPS ------------

func (s *FiberServer) StartFlagProcessingLoop(ctx context.Context) {
	interval := time.Duration(config.Current.ConfigServer.SubmitFlagCheckerTime) * time.Second
	ticker := time.NewTicker(interval)
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
			flags, err := s.db.GetUnsubmittedFlagCodeList(config.Current.ConfigServer.MaxFlagBatchSize)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to get unsubmitted flags")
				continue
			}
			if len(flags) == 0 {
				logger.Log.Debug().Msg("No flags to submit")
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
	maxBatch := int(config.Current.ConfigServer.MaxFlagBatchSize)
	groups := newFlagGroups(maxBatch)
	defer groups.reset()

	for _, resp := range flags {
		switch resp.Status {
		case models.StatusAccepted:
			groups.accepted = append(groups.accepted, resp.Flag)
		case models.StatusDenied:
			groups.denied = append(groups.denied, resp.Flag)
		default:
			groups.errored = append(groups.errored, resp.Flag)
		}
	}

	var wg sync.WaitGroup
	update := func(flags []string, status string) {
		defer wg.Done()
		if len(flags) == 0 {
			return
		}
		if err := s.db.UpdateFlagsStatus(flags, status); err != nil {
			logger.Log.Error().
				Strs("flags", flags).
				Err(err).
				Msgf("Failed to update flags with status %s", status)
		}
	}

	wg.Add(3)
	go update(groups.accepted, models.StatusAccepted)
	go update(groups.denied, models.StatusDenied)
	go update(groups.errored, models.StatusError)
	wg.Wait()

	total := len(groups.accepted) + len(groups.denied) + len(groups.errored)
	logger.Log.Info().
		Int("accepted", len(groups.accepted)).
		Int("denied", len(groups.denied)).
		Int("errored", len(groups.errored)).
		Int("total", total).
		Msg("Flags update summary")
}
