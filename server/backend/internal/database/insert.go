package database

import (
	"context"
	"time"

	"github.com/ByteTheCookies/backend/internal/models"
)

func (s *service) AddFlags(flags []models.Flag) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "INSERT INTO flags (id, flag_code, service_name,service_port, submit_time, response_time, status, team_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, flag := range flags {
		_, err := stmt.ExecContext(ctx, flag.ID, flag.FlagCode, flag.ServiceName, flag.ServicePort, flag.SubmitTime, flag.ResponseTime, flag.Status, flag.TeamID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) AddFlag(flag models.Flag) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "INSERT INTO flags (id, flag_code, service_name,service_port, submit_time, response_time, status, team_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, flag.ID, flag.FlagCode, flag.ServiceName, flag.ServicePort, flag.SubmitTime, flag.ResponseTime, flag.Status, flag.TeamID)
	if err != nil {
		return err
	}

	return nil
}
