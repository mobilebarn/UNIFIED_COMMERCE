package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/lib/pq"
)

// PostgresDB wraps a GORM database connection with additional utilities
type PostgresDB struct {
	DB     *gorm.DB
	SqlDB  *sql.DB
	Config *PostgresConfig
}

// PostgresConfig holds PostgreSQL configuration
type PostgresConfig struct {
	DatabaseURL     string // Full database connection URL (for cloud providers)
	Host            string
	Port            string
	User            string
	Password        string
	DatabaseName    string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	LogLevel        logger.LogLevel
}

// NewPostgresConnection creates a new PostgreSQL connection using GORM
func NewPostgresConnection(config *PostgresConfig) (*PostgresDB, error) {
	var dsn string

	// Check if we have a full database URL (for cloud providers like Render)
	if config.DatabaseURL != "" {
		dsn = config.DatabaseURL
	} else {
		// Build DSN from individual components
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host,
			config.Port,
			config.User,
			config.Password,
			config.DatabaseName,
			config.SSLMode,
		)
	}

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return &PostgresDB{
		DB:     db,
		SqlDB:  sqlDB,
		Config: config,
	}, nil
}

// Close closes the database connection
func (pdb *PostgresDB) Close() error {
	return pdb.SqlDB.Close()
}

// Health checks the database connection health
func (pdb *PostgresDB) Health(ctx context.Context) error {
	return pdb.SqlDB.PingContext(ctx)
}

// Migrate runs database migrations for the given models
func (pdb *PostgresDB) Migrate(models ...interface{}) error {
	return pdb.DB.AutoMigrate(models...)
}

// Transaction executes a function within a database transaction
func (pdb *PostgresDB) Transaction(fn func(*gorm.DB) error) error {
	return pdb.DB.Transaction(fn)
}

// DefaultPostgresConfig returns a default PostgreSQL configuration
func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:            "localhost",
		Port:            "5432",
		User:            "postgres",
		Password:        "postgres",
		DatabaseName:    "unified_commerce",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		LogLevel:        logger.Info,
	}
}

// NewPostgresConfigFromEnv creates PostgreSQL config from a unified config
func NewPostgresConfigFromEnv(host, port, user, password, dbName string) *PostgresConfig {
	config := DefaultPostgresConfig()
	config.Host = host
	config.Port = port
	config.User = user
	config.Password = password
	config.DatabaseName = dbName

	return config
}
