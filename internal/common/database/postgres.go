package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// NewPostgresConnectionWithConfig creates a new connection to the Postgres database
// using the provided configuration
func NewPostgresConnectionWithConfig(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.GetDBConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// Configure the connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
