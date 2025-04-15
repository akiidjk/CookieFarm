package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "embed"

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
	GetFlags() ([]models.Flag, error)
	GetAllFlags() ([]models.Flag, error)
	GetFlagsCode() ([]string, error)
	GetAllFlagsCode() ([]string, error)
	UpdateFlagStatus(flag_code string, status string) error
	UpdateFlagsStatus(flags []string, status string) error
	InitDB() error
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	dburl      = utils.GetEnv("DB_URL", "internal/database/cookiefarm.db")
	dbInstance *service
)

func (s *service) InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	logger.Info("Initializing database")
	_, err := s.db.ExecContext(ctx, sqlSchema)
	if err != nil {
		return err
	}
	logger.Info("Database initialized")
	return nil
}

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		logger.Fatal("Failed to open database connection %v", err)
	}

	dbInstance = &service{
		db: db,
	}

	if err := dbInstance.InitDB(); err != nil {
		logger.Fatal("InitDB failed: %v", err)
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
		logger.Fatal("Failed to ping database %v", err) // Log the error and terminate the program
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
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) Close() error {
	logger.Info("Disconnected from database: %s", dburl)
	return s.db.Close()
}
