package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/models"
)

func AddFlags(flags []models.Flag) error {
	if len(flags) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tx, _ := DB.BeginTx(ctx, nil)
	defer tx.Rollback()

	const maxParams = 50
	perRow := 8
	maxRows := maxParams / perRow
	for i := 0; i < len(flags); i += maxRows {
		end := i + maxRows
		if end > len(flags) {
			end = len(flags)
		}
		batch := flags[i:end]

		parts := make([]string, len(batch))
		args := make([]any, 0, len(batch)*perRow)
		for j, f := range batch {
			parts[j] = "(?, ?, ?, ?, ?, ?, ?)"
			args = append(args,
				f.FlagCode, f.ServiceName, f.ServicePort,
				f.SubmitTime, f.ResponseTime, f.Status, f.TeamID,
			)
		}
		query := "INSERT INTO flags(flag_code,service_name,service_port,submit_time,response_time,status,team_id) VALUES " +
			strings.Join(parts, ",")

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert: %w", err)
		}
	}

	flags = nil

	return tx.Commit()
}

func AddFlag(flag models.Flag) error {
	return AddFlags([]models.Flag{flag})
}
