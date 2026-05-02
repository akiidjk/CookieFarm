package database

import (
	"context"
	"database/sql"
	"strings"
)

type Store struct {
	db      *sql.DB
	Queries *Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) WithTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	if err := fn(q); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// INSERT OR IGNORE INTO flags(
// flag_code, service_name, port_service,
// submit_time, response_time, status,
// team_id, msg, username, exploit_name
// ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

func (s *Store) BulkInsertFlags(ctx context.Context, rows []Flag) error {
	if len(rows) == 0 {
		return nil
	}

	const maxParams = 900 // SQLite has a default limit of 999 parameters per statement, but we reserve some for safety
	paramsPerRow := 2
	maxRowsPerBatch := maxParams / paramsPerRow

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for start := 0; start < len(rows); start += maxRowsPerBatch {
		end := min(start+maxRowsPerBatch, len(rows))

		batch := rows[start:end]

		var (
			sb   strings.Builder
			args []any
		)

		sb.WriteString("INSERT OR IGNORE INTO flags(flag_code, service_name, port_service,submit_time, response_time, status, team_id, msg, username, exploit_name) VALUES ")

		for i, r := range batch {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			args = append(args,
				r.FlagCode, r.ServiceName, r.PortService,
				r.SubmitTime, r.ResponseTime, r.Status,
				r.TeamID, r.Msg, r.Username, r.ExploitName,
			)
		}
		sb.WriteString(";")

		stmt, err := tx.PrepareContext(ctx, sb.String())
		if err != nil {
			return err
		}

		defer stmt.Close()
		if _, err = stmt.ExecContext(ctx, args...); err != nil {
			return err
		}

	}

	return tx.Commit()
}
