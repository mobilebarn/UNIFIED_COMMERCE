package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Environment string
	Database    DatabaseConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig
}

// PostgreSQLConfig holds PostgreSQL configuration
type PostgreSQLConfig struct {
	URL string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Database: DatabaseConfig{
			PostgreSQL: PostgreSQLConfig{
				URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/unified_commerce?sslmode=disable"),
			},
		},
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
