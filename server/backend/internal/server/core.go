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
	logger.Info("Starting loop")

	var err error
	config.Submit, err = protocols.LoadProtocol("cc_http")

	if err != nil {
		logger.Error("LoadProtocol error: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Loop terminated")
			return
		case <-ticker.C:
			flags, err := s.db.GetUnsubmittedFlagsCode(int(config.Current.ConfigServer.MaxFlagBatchSize))
			if err != nil {
				logger.Error("GetUnsubmittedFlagsCode error: %v", err)
				continue
			}

			if len(flags) == 0 {
				logger.Debug("No flags to submit")
				continue
			}

			res, err := config.Submit(config.Current.ConfigServer.HostFlagchecker, config.Current.ConfigServer.TeamToken, flags)
			if err != nil {
				logger.Error("Submit error: %v", err)
				continue
			}

			s.UpdateFlags(res)
		}
	}
}

func (s *FiberServer) UpdateFlags(flags []models.ResponseProtocol) {
	flagsAccepted := make([]string, 0)
	flagsDenied := make([]string, 0)
	flagsError := make([]string, 0)
	for key, value := range flags {
		switch value.Status {
		case "ACCEPTED":
			flagsAccepted = append(flagsAccepted, value.Flag)
			logger.Debug("Flag %d accepted = %v", key, value.Flag)
		case "DENIED":
			flagsDenied = append(flagsDenied, value.Flag)
			logger.Debug("Flag %d denied = %v", key, value.Flag)
		default:
			flagsError = append(flagsError, value.Flag)
			logger.Debug("Flag %d error = %v", key, value.Flag)
		}
	}
	if len(flagsAccepted) > 0 {
		s.db.UpdateFlagsStatus(flagsAccepted, "ACCEPTED")
	}
	if len(flagsDenied) > 0 {
		s.db.UpdateFlagsStatus(flagsDenied, "DENIED")
	}
	if len(flagsError) > 0 {
		s.db.UpdateFlagsStatus(flagsError, "ERROR")
	}
	logger.Info("Flags updated. Accepted: %d, Denied: %d, Error: %d", len(flagsAccepted), len(flagsDenied), len(flagsError))

}
