package config

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all the configuration for the shopping-experience microservice
type Config struct {
	// Server configuration
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration

	// Database configuration
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
	DbSslMode  string

	// External services
	ProductCatalogServiceURL string
}

// LoadConfig loads the configuration from environment variables with appropriate prefixes
func LoadConfig() (*Config, error) {
	// 1. Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, loading environment variables from the system")
	}

	// 2. Set default values for configuration
	viper.SetDefault("SHOPPING_EXPERIENCE_HOST", "0.0.0.0")
	viper.SetDefault("SHOPPING_EXPERIENCE_PORT", 8001)
	viper.SetDefault("SHOPPING_EXPERIENCE_READ_TIMEOUT", "10s")
	viper.SetDefault("SHOPPING_EXPERIENCE_WRITE_TIMEOUT", "10s")
	viper.SetDefault("SHOPPING_EXPERIENCE_IDLE_TIMEOUT", "120s")
	viper.SetDefault("SHOPPING_EXPERIENCE_SHUTDOWN_TIMEOUT", "15s")
	viper.SetDefault("SHOPPING_EXPERIENCE_DB_HOST", "localhost")
	viper.SetDefault("SHOPPING_EXPERIENCE_DB_PORT", 5432)
	viper.SetDefault("SHOPPING_EXPERIENCE_DB_USER", "shopping_experience")
	viper.SetDefault("SHOPPING_EXPERIENCE_DB_PASS", "shopping_experience")
	viper.SetDefault("SHOPPING_EXPERIENCE_DB_NAME", "shopping_experience")
	viper.SetDefault("SHOPPING_EXPERIENCE_DB_SSLMODE", "disable")
	viper.SetDefault("PRODUCT_CATALOG_SERVICE_URL", "http://localhost:8000")

	// 3. Get configuration values from environment variables
	viper.AutomaticEnv()

	// Parse durations from string values
	readTimeout, err := time.ParseDuration(viper.GetString("SHOPPING_EXPERIENCE_READ_TIMEOUT"))
	if err != nil {
		readTimeout = 10 * time.Second
	}

	writeTimeout, err := time.ParseDuration(viper.GetString("SHOPPING_EXPERIENCE_WRITE_TIMEOUT"))
	if err != nil {
		writeTimeout = 10 * time.Second
	}

	idleTimeout, err := time.ParseDuration(viper.GetString("SHOPPING_EXPERIENCE_IDLE_TIMEOUT"))
	if err != nil {
		idleTimeout = 120 * time.Second
	}

	shutdownTimeout, err := time.ParseDuration(viper.GetString("SHOPPING_EXPERIENCE_SHUTDOWN_TIMEOUT"))
	if err != nil {
		shutdownTimeout = 15 * time.Second
	}

	config := &Config{
		Host:                     viper.GetString("SHOPPING_EXPERIENCE_HOST"),
		Port:                     viper.GetInt("SHOPPING_EXPERIENCE_PORT"),
		ReadTimeout:              readTimeout,
		WriteTimeout:             writeTimeout,
		IdleTimeout:              idleTimeout,
		ShutdownTimeout:          shutdownTimeout,
		DbHost:                   viper.GetString("SHOPPING_EXPERIENCE_DB_HOST"),
		DbPort:                   viper.GetInt("SHOPPING_EXPERIENCE_DB_PORT"),
		DbUser:                   viper.GetString("SHOPPING_EXPERIENCE_DB_USER"),
		DbPassword:               viper.GetString("SHOPPING_EXPERIENCE_DB_PASS"),
		DbName:                   viper.GetString("SHOPPING_EXPERIENCE_DB_NAME"),
		DbSslMode:                viper.GetString("SHOPPING_EXPERIENCE_DB_SSLMODE"),
		ProductCatalogServiceURL: viper.GetString("PRODUCT_CATALOG_SERVICE_URL"),
	}

	return config, nil
}

// GetDBConnectionString returns a formatted database connection string
func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DbHost, c.DbPort, c.DbUser, c.DbPassword, c.DbName, c.DbSslMode)
}
