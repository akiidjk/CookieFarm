package database

import (
	"context"
	"time"

	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
)

const insertFlagQuery = `
	INSERT INTO flags
	(id, flag_code, service_name, service_port, submit_time, response_time, status, team_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

func (s *service) AddFlags(flags []models.Flag) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, insertFlagQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to prepare insert statement for flags")
		return err
	}
	defer stmt.Close()

	for _, flag := range flags {
		_, err := stmt.ExecContext(ctx,
			flag.ID,
			flag.FlagCode,
			flag.ServiceName,
			flag.ServicePort,
			flag.SubmitTime,
			flag.ResponseTime,
			flag.Status,
			flag.TeamID,
		)
		if err != nil {
			logger.Log.Error().Err(err).Str("flag_id", flag.ID).Msg("Failed to insert flag")
			return err
		}
	}

	logger.Log.Debug().Int("inserted", len(flags)).Msg("Flags inserted successfully")
	return nil
}

func (s *service) AddFlag(flag models.Flag) error {
	return s.AddFlags([]models.Flag{flag})
}
