package database

import (
	"context"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
)

const (
	baseFlagQuery         = `SELECT flag_code, service_name,service_port, submit_time, response_time, status, team_id FROM flags`
	queryAllFlags         = baseFlagQuery + " ORDER BY submit_time DESC"
	queryFirstNFlags      = baseFlagQuery + " ORDER BY submit_time DESC LIMIT ?"
	queryUnsubmittedFlags = baseFlagQuery + " WHERE status = 'UNSUBMITTED' ORDER BY submit_time ASC LIMIT ?"
	queryPagedFlags       = baseFlagQuery + " ORDER BY submit_time DESC LIMIT ? OFFSET ?"

	baseFlagCodeQuery         = `SELECT flag_code FROM flags`
	queryAllFlagCodes         = baseFlagCodeQuery
	queryFirstNFlagCodes      = baseFlagCodeQuery + " LIMIT ?"
	queryUnsubmittedFlagCodes = baseFlagCodeQuery + " WHERE status = 'UNSUBMITTED' LIMIT ?"
	queryPagedFlagCodes       = baseFlagCodeQuery + " LIMIT ? OFFSET ?"
)

// --------- Flag Structs ---------

func GetAllFlags() ([]models.Flag, error) {
	return queryFlags(queryAllFlags)
}

func GetUnsubmittedFlags(limit int) ([]models.Flag, error) {
	return queryFlags(queryUnsubmittedFlags, limit)
}

func GetFirstNFlags(limit int) ([]models.Flag, error) {
	return queryFlags(queryFirstNFlags, limit)
}

func GetPagedFlags(limit, offset int) ([]models.Flag, error) {
	return queryFlags(queryPagedFlags, limit, offset)
}

// --------- Flag Code Only ---------

func GetAllFlagCodeList() ([]string, error) {
	return queryFlagCodes(queryAllFlagCodes)
}

func GetUnsubmittedFlagCodeList(limit uint16) ([]string, error) {
	return queryFlagCodes(queryUnsubmittedFlagCodes, limit)
}

func GetFirstNFlagCodeList(limit int) ([]string, error) {
	return queryFlagCodes(queryFirstNFlagCodes, limit)
}

func GetPagedFlagCodeList(limit, offset int) ([]string, error) {
	return queryFlagCodes(queryPagedFlagCodes, limit, offset)
}

// --------- Shared query logic ---------

func queryFlags(query string, args ...any) ([]models.Flag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stmt, err := DB.PrepareContext(ctx, query)
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
	flagPtr := new(models.Flag)
	for rows.Next() {
		if err := rows.Scan(
			&flagPtr.FlagCode, &flagPtr.ServiceName, &flagPtr.ServicePort,
			&flagPtr.SubmitTime, &flagPtr.ResponseTime, &flagPtr.Status,
			&flagPtr.TeamID,
		); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan row in queryFlags")
			return nil, err
		}
		flags = append(flags, *flagPtr)
	}

	return flags, nil
}

func queryFlagCodes(query string, args ...any) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stmt, err := DB.PrepareContext(ctx, query)
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
	codePtr := new(string)
	for rows.Next() {
		if err := rows.Scan(codePtr); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan row in queryFlagCodes")
			return nil, err
		}
		codes = append(codes, *codePtr)
	}

	return codes, nil
}
