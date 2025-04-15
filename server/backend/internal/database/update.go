package database

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (s *service) UpdateFlagStatus(flagCode string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "UPDATE flags SET status = ? WHERE flag_code = ?"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, status, flagCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateFlagsStatus(flagsCode []string, status string) error {
	if len(flagsCode) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	placeholders := strings.Repeat("?,", len(flagsCode))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf("UPDATE flags SET status = ? WHERE flag_code IN (%s)", placeholders)

	args := make([]interface{}, 0, len(flagsCode)+1)
	args = append(args, status)
	for _, code := range flagsCode {
		args = append(args, code)
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	return err
}
