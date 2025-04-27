package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/logger"
)

const (
	updateFlagStatusQuery  = `UPDATE flags SET status = ?, response_time = ? WHERE flag_code = ?`
	updateFlagsStatusQuery = `UPDATE flags SET status = ?, response_time = ? WHERE flag_code IN (%s)`
)

func (s *service) UpdateFlagStatus(flagCode string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, updateFlagStatusQuery)
	if err != nil {
		logger.Log.Error().Err(err).Str("flag_code", flagCode).Msg("Failed to prepare UpdateFlagStatus statement")
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, status, uint64(time.Now().Unix()), flagCode)
	if err != nil {
		logger.Log.Error().Err(err).Str("flag_code", flagCode).Msg("Failed to execute UpdateFlagStatus")
		return err
	}

	logger.Log.Debug().Str("flag_code", flagCode).Str("status", status).Msg("Updated flag status")
	return nil
}

func (s *service) UpdateFlagsStatus(flagCodes []string, status string) error {
	if len(flagCodes) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	placeholders := strings.Repeat("?,", len(flagCodes))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(updateFlagsStatusQuery, placeholders)

	args := make([]any, 0, len(flagCodes)+2)
	args = append(args, status, time.Now().Unix())
	for _, code := range flagCodes {
		args = append(args, code)
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		logger.Log.Error().Err(err).Int("count", len(flagCodes)).Msg("Failed to prepare UpdateFlagsStatus statement")
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		logger.Log.Error().Err(err).Int("count", len(flagCodes)).Msg("Failed to execute UpdateFlagsStatus")
		return err
	}

	logger.Log.Debug().Int("count", len(flagCodes)).Str("status", status).Msg("Updated statuses for multiple flags")
	return nil
}
