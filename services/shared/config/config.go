package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration values for a microservice
type Config struct {
	// Service configuration
	ServiceName string
	ServicePort string
	Environment string

	// Database configuration
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	// MongoDB configuration (for services that use MongoDB)
	MongoURL      string
	MongoDatabase string
	MongoUser     string
	MongoPassword string

	// Redis configuration
	RedisURL      string
	RedisPassword string
	RedisDB       int

	// JWT configuration
	JWTSecret     string
	JWTExpiration int

	// Observability
	JaegerEndpoint string
	LogLevel       string

	// GraphQL Gateway
	GatewayURL string

	// External services
	ElasticsearchURL string
	KafkaBrokers     []string

	// Development flags
	DebugMode bool
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig(serviceName string) (*Config, error) {
	// Load .env file if it exists (for local development)
	// Try to find .env file in current directory or parent directories
	envFilePath := findEnvFile()
	if envFilePath != "" {
		loadEnvWithoutBOM(envFilePath)
	}

	// Support both DB_* and DATABASE_* prefixes; sanitize service name for default DB
	sanitizeServiceDBName := func(name string) string {
		// Postgres identifiers can't easily use hyphens unless quoted; replace with underscore
		return strings.ReplaceAll(name, "-", "_")
	}

	getEnvAny := func(keys []string, def string) string {
		for _, k := range keys {
			if v := os.Getenv(k); v != "" {
				return v
			}
		}
		return def
	}

	config := &Config{
		ServiceName: serviceName,
		ServicePort: getEnv("SERVICE_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		DatabaseHost:     getEnvAny([]string{"DB_HOST", "DATABASE_HOST"}, "localhost"),
		DatabasePort:     getEnvAny([]string{"DB_PORT", "DATABASE_PORT"}, "5432"),
		DatabaseUser:     getEnvAny([]string{"DB_USER", "DATABASE_USER"}, "postgres"),
		DatabasePassword: getEnvAny([]string{"DB_PASSWORD", "DATABASE_PASSWORD"}, "postgres"),
		DatabaseName:     getEnvAny([]string{"DB_NAME", "DATABASE_NAME"}, sanitizeServiceDBName(serviceName)+"_service"),

		// MongoDB
		MongoURL:      getEnv("MONGO_URL", "mongodb://mongodb:mongodb@localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "product_catalog"),
		MongoUser:     getEnv("MONGO_USER", "catalog_user"),
		MongoPassword: getEnv("MONGO_PASSWORD", "catalog_pass"),

		// Redis
		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", "redis"),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "unified-commerce-jwt-secret-key"),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 3600), // 1 hour

		// Observability
		JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),

		// Gateway
		GatewayURL: getEnv("GATEWAY_URL", "http://localhost:4000"),

		// External services
		ElasticsearchURL: getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		KafkaBrokers:     []string{getEnv("KAFKA_BROKERS", "localhost:9092")},

		// Development
		DebugMode: getEnvAsBool("DEBUG_MODE", true),
	}

	// Build database URL
	// Check if DATABASE_URL is provided directly (for Render, Railway, etc.)
	if envDatabaseURL := os.Getenv("DATABASE_URL"); envDatabaseURL != "" {
		config.DatabaseURL = envDatabaseURL
	} else {
		config.DatabaseURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseUser,
			config.DatabasePassword,
			config.DatabaseName,
		)
	}

	return config, nil
}

// findEnvFile attempts to find .env file in current directory or parent directories
func findEnvFile() string {
	// For payment service, look in the specific service directory first
	paths := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
		"services/payment/.env", // Add specific path for payment service
		"../services/payment/.env",
		"../../services/payment/.env",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Found .env file at %s\n", path)
			return path
		}
	}

	return ""
}

// loadEnvWithoutBOM loads .env file and handles BOM if present
func loadEnvWithoutBOM(filepath string) {
	// Read the entire file
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading .env file: %v\n", err)
		return
	}

	// Remove BOM if present (UTF-8 BOM is 0xEF,0xBB,0xBF)
	// Handle multiple BOM sequences if present
	for len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		// BOM found, remove it
		data = data[3:]
		fmt.Println("BOM found and removed from .env file")
	}

	// Parse the env file without BOM
	envMap, err := godotenv.Unmarshal(string(data))
	if err != nil {
		fmt.Printf("Error parsing .env file: %v\n", err)
		return
	}

	// Set environment variables
	for key, value := range envMap {
		os.Setenv(key, value)
	}

	fmt.Printf("Successfully loaded .env file, SERVICE_PORT: '%s'\n", os.Getenv("SERVICE_PORT"))
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
