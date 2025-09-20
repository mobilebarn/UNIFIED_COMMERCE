package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// containsCredentials checks if the MongoDB URI contains username and password
func containsCredentials(uri string) bool {
	return strings.Contains(uri, "@")
}

// MongoConnection represents a MongoDB connection

// MongoDB wraps a MongoDB client with additional utilities
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	Config   *MongoConfig
}

// MongoConfig holds MongoDB configuration
type MongoConfig struct {
	URI            string
	DatabaseName   string
	Username       string
	Password       string
	MaxPoolSize    uint64
	MinPoolSize    uint64
	ConnectTimeout time.Duration
	ServerTimeout  time.Duration
}

// NewMongoConnection creates a new MongoDB connection
func NewMongoConnection(config *MongoConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	// Debug logging
	fmt.Printf("MongoDB Connection Debug:\n")
	fmt.Printf("  URI: %s\n", config.URI[:min(60, len(config.URI))])
	fmt.Printf("  Database: %s\n", config.DatabaseName)
	fmt.Printf("  Username: %s\n", config.Username)
	fmt.Printf("  Has credentials in URI: %v\n", containsCredentials(config.URI))

	// Configure client options - use URI directly if it contains credentials
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize).
		SetServerSelectionTimeout(config.ServerTimeout)

	// Only add separate auth if URI doesn't contain credentials and we have username/password
	if config.Username != "" && config.Password != "" && !containsCredentials(config.URI) {
		fmt.Printf("  Adding separate authentication credentials\n")
		credential := options.Credential{
			Username:   config.Username,
			Password:   config.Password,
			AuthSource: "admin", // MongoDB cloud services typically require auth against admin db
		}
		clientOptions.SetAuth(credential)
	} else {
		fmt.Printf("  Using URI-embedded credentials or no auth\n")
	}

	// Connect to MongoDB
	fmt.Printf("  Attempting connection...\n")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("  Connection failed: %v\n", err)
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	fmt.Printf("  Testing connection with ping...\n")
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Printf("  Ping failed: %v\n", err)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Printf("  MongoDB connection successful!\n")

	database := client.Database(config.DatabaseName)

	return &MongoDB{
		Client:   client,
		Database: database,
		Config:   config,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// Health checks the MongoDB connection health
func (m *MongoDB) Health(ctx context.Context) error {
	return m.Client.Ping(ctx, readpref.Primary())
}

// Collection returns a MongoDB collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

// StartTransaction starts a new MongoDB transaction session
func (m *MongoDB) StartTransaction(ctx context.Context) (mongo.Session, error) {
	session, err := m.Client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	if err := session.StartTransaction(); err != nil {
		session.EndSession(ctx)
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return session, nil
}

// WithTransaction executes a function within a MongoDB transaction
func (m *MongoDB) WithTransaction(ctx context.Context, fn func(mongo.SessionContext) error) error {
	session, err := m.Client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	callback := func(sc mongo.SessionContext) (interface{}, error) {
		return nil, fn(sc)
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

// DefaultMongoConfig returns a default MongoDB configuration
func DefaultMongoConfig() *MongoConfig {
	return &MongoConfig{
		URI:            "mongodb://localhost:27017",
		DatabaseName:   "product_catalog",
		Username:       "",
		Password:       "",
		MaxPoolSize:    100,
		MinPoolSize:    5,
		ConnectTimeout: 10 * time.Second,
		ServerTimeout:  5 * time.Second,
	}
}

// NewMongoConfigFromEnv creates MongoDB config from environment variables
func NewMongoConfigFromEnv(uri, database, username, password string) *MongoConfig {
	config := DefaultMongoConfig()
	config.URI = uri
	config.DatabaseName = database
	config.Username = username
	config.Password = password

	return config
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
