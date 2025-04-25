package database

import (
	"context"
	"strings"
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
	if len(flags) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := `
		INSERT INTO flags
		(id, flag_code, service_name, service_port, submit_time, response_time, status, team_id)
		VALUES `

	valueStrings := make([]string, 0, len(flags))
	valueArgs := make([]interface{}, 0, len(flags)*8)

	for _, flag := range flags {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			flag.ID,
			flag.FlagCode,
			flag.ServiceName,
			flag.ServicePort,
			flag.SubmitTime,
			flag.ResponseTime,
			flag.Status,
			flag.TeamID,
		)
	}

	query += strings.Join(valueStrings, ", ")

	_, err := s.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to batch insert flags")
		return err
	}

	logger.Log.Debug().Int("inserted", len(flags)).Msg("Flags inserted successfully")
	return nil
}

func (s *service) AddFlag(flag models.Flag) error {
	return s.AddFlags([]models.Flag{flag})
}
