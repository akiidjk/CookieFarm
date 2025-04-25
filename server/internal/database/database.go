package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/ByteTheCookies/backend/internal/models"
	"github.com/ByteTheCookies/backend/internal/utils"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var sqlSchema string

type Service interface {
	Health() map[string]string
	AddFlags(flags []models.Flag) error
	AddFlag(flag models.Flag) error
	GetUnsubmittedFlags(limit int) ([]models.Flag, error)
	GetAllFlags() ([]models.Flag, error)
	GetFirstNFlags(limit int) ([]models.Flag, error)
	GetPagedFlags(offset int, limit int) ([]models.Flag, error)
	GetPagedFlagCodeList(offset int, limit int) ([]string, error)
	GetUnsubmittedFlagCodeList(limit int) ([]string, error)
	GetAllFlagCodeList() ([]string, error)
	GetFirstNFlagCodeList(limit int) ([]string, error)
	UpdateFlagStatus(flag_code string, status string) error
	UpdateFlagsStatus(flags []string, status string) error
	FlagsNumber(ctx context.Context) int
	InitDB() error
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	dbPath     = utils.GetEnv("DB_URL", filepath.Join(utils.GetExecutableDir(), "cookiefarm.db"))
	dbInstance *service
)

func (s *service) InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	logger.Log.Info().Msg("Initializing database schema")
	_, err := s.db.ExecContext(ctx, sqlSchema)
	if err != nil {
		return err
	}

	logger.Log.Info().Msg("Database schema initialized successfully")
	return nil
}

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Log.Fatal().Err(err).Str("path", dbPath).Msg("Failed to open database")
	}

	dbInstance = &service{
		db: db,
	}

	if err := dbInstance.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database initialization failed")
	}

	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		logger.Log.Error().Err(err).Msg("Database ping failed")
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed. Consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime. Consider increasing it."
	}

	return stats
}

func (s *service) Close() error {
	logger.Log.Info().Str("path", dbPath).Msg("Disconnected from database")
	return s.db.Close()
}

func (s *service) FlagsNumber(ctx context.Context) int {

	stmt, err := s.db.PrepareContext(ctx, "SELECT COUNT(*) FROM flags")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to prepare statement")
		return 0
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to execute query")
		return 0
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan result")
			return 0
		}
	}

	logger.Log.Debug().Int("count", count).Msg("Flags number")

	return count
}
