package server

import (
	"context"
	"time"

	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/protocols"
)

func (s *FiberServer) StartLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Loop terminated")
			return
		case <-ticker.C:
			flags, err := s.db.GetFlagsCode()
			if err != nil {
				logger.Error("GetFlagsCode error: %v", err)
				continue
			}

			res, err := protocols.Submit(flags)
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
	for key, value := range flags {
		logger.Debug("flag %d = %v", key, value.Flag)
		if value.Status == "ACCEPTED" {
			logger.Info("Flag %d accepted", key)
			flagsAccepted = append(flagsAccepted, value.Flag)
		}
	}
	s.db.UpdateFlagsStatus(flagsAccepted, "ACCEPTED")
}
