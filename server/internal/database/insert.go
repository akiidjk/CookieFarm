package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieserver/internal/models"
)

// AddFlags adds a batch of flags to the database.
// It divides the flags into batches to insert in chunks, helping avoid hitting query parameter limits.
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
				f.FlagCode, f.ServiceName, f.PortService,
				f.SubmitTime, f.ResponseTime, f.Status, f.TeamID,
			)
		}
		query := "INSERT INTO flags(flag_code,service_name,port_service,submit_time,response_time,status,team_id) VALUES " +
			strings.Join(parts, ",")

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert: %w", err)
		}
	}

	flags = nil

	return tx.Commit()
}

// AddFlag adds a single flag to the database.
// It calls the AddFlags function to add the flag as a batch of size 1.
func AddFlag(flag models.Flag) error {
	return AddFlags([]models.Flag{flag})
}
