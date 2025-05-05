package server

import (
	"context"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/database"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/ByteTheCookies/cookieserver/protocols"
)

// ----------- FLAG GROUPS ------------

// FlagGroups represents a group of flags with different statuses (accepted, denied, errored).
type FlagGroups struct {
	accepted []string // Accepted flags
	denied   []string // Denied flags
	errored  []string // Errored flags
	cap      int      // Capacity of the groups
}

// newFlagGroups creates a new instance of FlagGroups with the specified capacity.
func newFlagGroups(cap int) *FlagGroups {
	return &FlagGroups{
		accepted: make([]string, 0, cap),
		denied:   make([]string, 0, cap),
		errored:  make([]string, 0, cap),
		cap:      cap,
	}
}

// reset resets the flag groups to their initial state.
func (g *FlagGroups) reset() {
	g.accepted = g.accepted[:0]
	g.denied = g.denied[:0]
	g.errored = g.errored[:0]
}

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
		select {
		case <-ctx.Done():
			logger.Log.Info().Msg("Flag processing loop terminated")
			return
		case <-ticker.C:
			flags, err := database.GetUnsubmittedFlagCodeList(config.Current.ConfigServer.MaxFlagBatchSize)
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
func UpdateFlags(flags []models.ResponseProtocol) {
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
		if err := database.UpdateFlagsStatus(flags, status); err != nil {
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

// LoadConfig loads the configuration from the given path.
func LoadConfig(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config.Current)
	if err != nil {
		return err
	}

	if config.Current.Configured != true {
		config.Current.Configured = false
	}

	if shutdownCancel != nil {
		shutdownCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	shutdownCancel = cancel

	go StartFlagProcessingLoop(ctx)

	return nil
}
