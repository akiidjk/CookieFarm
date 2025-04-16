package database

import (
	"context"
	"time"

	"github.com/ByteTheCookies/backend/internal/models"
)

func (s *service) GetAllFlags() ([]models.Flag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	flags := []models.Flag{}

	query := "SELECT id, flag_code, service_name, submit_time, status, team_id FROM flags"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flag models.Flag
		if err := rows.Scan(&flag.ID, &flag.FlagCode, &flag.ServiceName, &flag.SubmitTime, &flag.Status, &flag.TeamID); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (s *service) GetUnsubmittedFlags(limit int) ([]models.Flag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	flags := []models.Flag{}

	query := "SELECT id, flag_code, service_name, submit_time, status, team_id FROM flags WHERE status = 'UNSUBMITTED' LIMIT ?"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flag models.Flag
		if err := rows.Scan(&flag.ID, &flag.FlagCode, &flag.ServiceName, &flag.SubmitTime, &flag.Status, &flag.TeamID); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (s *service) GetAllFlagsCode() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	flags := []string{}

	query := "SELECT flag_code FROM flags"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flag string
		if err := rows.Scan(&flag); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (s *service) GetUnsubmittedFlagsCode(limit int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	flags := []string{}

	query := "SELECT flag_code FROM flags WHERE status = 'UNSUBMITTED' LIMIT ?"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flag string
		if err := rows.Scan(&flag); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}
