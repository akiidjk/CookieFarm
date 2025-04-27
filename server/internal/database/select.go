package database

import (
	"context"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
)

const (
	queryAllFlags             = `SELECT id, flag_code, service_name, submit_time, response_time, status, team_id FROM flags`
	queryAllFlagCodes         = `SELECT flag_code FROM flags`
	queryFirstNFlags          = queryAllFlags + ` LIMIT ?`
	queryFirstNFlagCodes      = queryAllFlagCodes + ` LIMIT ?`
	queryUnsubmittedFlags     = queryAllFlags + ` WHERE status = 'UNSUBMITTED' LIMIT ?`
	queryUnsubmittedFlagCodes = queryAllFlagCodes + ` WHERE status = 'UNSUBMITTED' LIMIT ?`
	queryPagedFlags           = queryAllFlags + ` LIMIT ? OFFSET ?`
	queryPagedFlagCodes       = queryAllFlagCodes + ` LIMIT ? OFFSET ?`
)

// --------- Flag Structs ---------

func (s *service) GetAllFlags() ([]models.Flag, error) {
	return s.queryFlags(queryAllFlags)
}

func (s *service) GetUnsubmittedFlags(limit int) ([]models.Flag, error) {
	return s.queryFlags(queryUnsubmittedFlags, limit)
}

func (s *service) GetFirstNFlags(limit int) ([]models.Flag, error) {
	return s.queryFlags(queryFirstNFlags, limit)
}

func (s *service) GetPagedFlags(limit, offset int) ([]models.Flag, error) {
	return s.queryFlags(queryPagedFlags, limit, offset)
}

// --------- Flag Code Only ---------

func (s *service) GetAllFlagCodeList() ([]string, error) {
	return s.queryFlagCodes(queryAllFlagCodes)
}

func (s *service) GetUnsubmittedFlagCodeList(limit uint16) ([]string, error) {
	return s.queryFlagCodes(queryUnsubmittedFlagCodes, limit)
}

func (s *service) GetFirstNFlagCodeList(limit int) ([]string, error) {
	return s.queryFlagCodes(queryFirstNFlagCodes, limit)
}

func (s *service) GetPagedFlagCodeList(limit, offset int) ([]string, error) {
	return s.queryFlagCodes(queryPagedFlagCodes, limit, offset)
}

// --------- Shared query logic ---------

func (s *service) queryFlags(query string, args ...any) ([]models.Flag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		logger.Log.Error().Err(err).Str("query", query).Msg("Failed to prepare queryFlags")
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.Log.Error().Err(err).Str("query", query).Msg("Failed to execute queryFlags")
		return nil, err
	}
	defer rows.Close()

	var flags []models.Flag
	for rows.Next() {
		var flag models.Flag
		if err := rows.Scan(&flag.ID, &flag.FlagCode, &flag.ServiceName, &flag.SubmitTime, &flag.ResponseTime, &flag.Status, &flag.TeamID); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan row in queryFlags")
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (s *service) queryFlagCodes(query string, args ...any) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		logger.Log.Error().Err(err).Str("query", query).Msg("Failed to prepare queryFlagCodes")
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.Log.Error().Err(err).Str("query", query).Msg("Failed to execute queryFlagCodes")
		return nil, err
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan row in queryFlagCodes")
			return nil, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}
