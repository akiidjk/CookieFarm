package server

import (
	"context"
	"time"

	"github.com/ByteTheCookies/backend/internal/config"
	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/protocols"
)

func (s *FiberServer) StartLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	logger.Info("Starting loop")
	for {
		select {
		case <-ctx.Done():
			logger.Info("Loop terminated")
			return
		case <-ticker.C:
			flags, err := s.db.GetUnsubmittedFlagsCode(config.MAX_FLAG_BATCH_SIZE)
			if err != nil {
				logger.Error("GetUnsubmittedFlagsCode error: %v", err)
				continue
			}

			res, err := protocols.Submit(config.HOST, config.TEAM_TOKEN, flags)
			if err != nil {
				logger.Error("Submit error: %v", err)
				continue
			}

			s.UpdateFlags(res)

			// logger.Debug("Submit results: %v", res)
		}
	}
}

func (s *FiberServer) UpdateFlags(flags []models.ResponseProtocol) {
	flagsAccepted := make([]string, 0)
	flagsDenied := make([]string, 0)
	flagsError := make([]string, 0)
	logger.Debug("Total flag submitted: %d", len(flags))
	for key, value := range flags {
		if value.Status == "ACCEPTED" {
			logger.Debug("Flag %d accepted = %v", key, value.Flag)
			flagsAccepted = append(flagsAccepted, value.Flag)
		} else if value.Status == "DENIED" {
			logger.Debug("Flag %d denied = %v", key, value.Flag)
			flagsDenied = append(flagsDenied, value.Flag)
		} else {
			logger.Debug("Flag %d error = %v", key, value.Flag)
			flagsError = append(flagsError, value.Flag)
		}
	}
	s.db.UpdateFlagsStatus(flagsAccepted, "ACCEPTED")
	s.db.UpdateFlagsStatus(flagsDenied, "DENIED")
	s.db.UpdateFlagsStatus(flagsError, "ERROR")
}
