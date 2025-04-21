package database

import (
	"context"
	"time"

	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
)

const (
	selectAllFlags             = `SELECT id, flag_code, service_name, submit_time, response_time, status, team_id FROM flags`
	selectUnsubmittedFlags     = selectAllFlags + ` WHERE status = 'UNSUBMITTED' LIMIT ?`
	selectAllFlagsCode         = `SELECT flag_code FROM flags`
	selectUnsubmittedFlagsCode = selectAllFlagsCode + ` WHERE status = 'UNSUBMITTED' LIMIT ?`
)

func (s *service) GetAllFlags() ([]models.Flag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, selectAllFlags)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to prepare GetAllFlags query")
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to execute GetAllFlags query")
		return nil, err
	}
	defer rows.Close()

	var flags []models.Flag
	for rows.Next() {
		var flag models.Flag
		if err := rows.Scan(&flag.ID, &flag.FlagCode, &flag.ServiceName, &flag.SubmitTime, &flag.ResponseTime, &flag.Status, &flag.TeamID); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan row in GetAllFlags")
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (s *service) GetUnsubmittedFlags(limit int) ([]models.Flag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, selectUnsubmittedFlags)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to prepare GetUnsubmittedFlags query")
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, limit)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to execute GetUnsubmittedFlags query")
		return nil, err
	}
	defer rows.Close()

	var flags []models.Flag
	for rows.Next() {
		var flag models.Flag
		if err := rows.Scan(&flag.ID, &flag.FlagCode, &flag.ServiceName, &flag.SubmitTime, &flag.ResponseTime, &flag.Status, &flag.TeamID); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan row in GetUnsubmittedFlags")
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (s *service) GetAllFlagsCode() ([]string, error) {
	return s.getFlagCodes(selectAllFlagsCode)
}

func (s *service) GetUnsubmittedFlagsCode(limit int) ([]string, error) {
	return s.getFlagCodes(selectUnsubmittedFlagsCode, limit)
}

func (s *service) getFlagCodes(query string, args ...any) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to prepare getFlagCodes query")
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to execute getFlagCodes query")
		return nil, err
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan row in getFlagCodes")
			return nil, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}
