package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

const (
	defaultHost     = "localhost"
	defaultPort     = 5432
	defaultUser     = "postgres"
	defaultPassword = "postgres"
	defaultDBName   = "shopping_experience"
	defaultSSLMode  = "disable"
)

// NewPostgresConnection creates a new connection to the Postgres database
func NewPostgresConnection() (*sql.DB, error) {
	host := getEnv("DB_HOST", defaultHost)
	port := getEnvAsInt("DB_PORT", defaultPort)
	user := getEnv("DB_USER", defaultUser)
	password := getEnv("DB_PASSWORD", defaultPassword)
	dbname := getEnv("DB_NAME", defaultDBName)
	sslmode := getEnv("DB_SSLMODE", defaultSSLMode)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

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

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := fmt.Sscanf(value, "%d", new(int)); err == nil && intValue > 0 {
			return intValue
		}
	}
	return defaultValue
}
