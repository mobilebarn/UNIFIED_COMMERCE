package database

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps a Redis client with additional utilities
type RedisClient struct {
	Client *redis.Client
	Config *RedisConfig
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Address      string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewRedisConnection creates a new Redis connection
func NewRedisConnection(config *RedisConfig) (*RedisClient, error) {
	options := &redis.Options{
		Addr:         config.Address,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	client := redis.NewClient(options)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		Client: client,
		Config: config,
	}, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// Health checks the Redis connection health
func (r *RedisClient) Health(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// Set stores a key-value pair with optional expiration
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Delete removes a key
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, keys...).Result()
}

// Expire sets an expiration time for a key
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

// HSet sets a hash field
func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.HSet(ctx, key, values...).Err()
}

// HGet gets a hash field value
func (r *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

// HGetAll gets all hash fields and values
func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

// HDel deletes hash fields
func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	return r.Client.HDel(ctx, key, fields...).Err()
}

// Increment increments a key's value
func (r *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// IncrementBy increments a key's value by a specific amount
func (r *RedisClient) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.Client.IncrBy(ctx, key, value).Result()
}

// SAdd adds members to a set
func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.Client.SAdd(ctx, key, members...).Err()
}

// SMembers gets all members of a set
func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(ctx, key).Result()
}

// SRem removes members from a set
func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.Client.SRem(ctx, key, members...).Err()
}

// ZAdd adds members to a sorted set
func (r *RedisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return r.Client.ZAdd(ctx, key, members...).Err()
}

// ZRange gets a range of members from a sorted set
func (r *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.ZRange(ctx, key, start, stop).Result()
}

// DefaultRedisConfig returns a default Redis configuration
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Address:      "localhost:6379",
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 3,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// NewRedisConfigFromEnv creates Redis config from environment variables
// If redisURL is a full Redis URL (redis://...), it will be parsed
// Otherwise, it will be treated as a simple address
func NewRedisConfigFromEnv(redisURL, password string, db int) *RedisConfig {
	config := DefaultRedisConfig()
	
	// Parse Redis URL if it starts with redis://
	if strings.HasPrefix(redisURL, "redis://") {
		// Parse the URL to extract components
		if parsed, err := parseRedisURL(redisURL); err == nil {
			config.Address = parsed.Address
			config.Password = parsed.Password
			if parsed.DB >= 0 {
				config.DB = parsed.DB
			}
		} else {
			// Fallback to using the URL as address (for backwards compatibility)
			config.Address = strings.TrimPrefix(redisURL, "redis://")
		}
	} else {
		// Simple address format
		config.Address = redisURL
		config.Password = password
		config.DB = db
	}

	return config
}

// parsedRedisConfig holds parsed Redis URL components
type parsedRedisConfig struct {
	Address  string
	Password string
	DB       int
}

// parseRedisURL parses a Redis URL and extracts connection components
// Format: redis://[:password@]host:port[/db]
func parseRedisURL(redisURL string) (*parsedRedisConfig, error) {
	parsed, err := url.Parse(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	config := &parsedRedisConfig{
		Address: parsed.Host,
		DB:      0, // default
	}

	// Extract password from user info
	if parsed.User != nil {
		if password, hasPassword := parsed.User.Password(); hasPassword {
			config.Password = password
		}
	}

	// Extract database number from path
	if parsed.Path != "" && parsed.Path != "/" {
		dbStr := strings.TrimPrefix(parsed.Path, "/")
		if db, err := strconv.Atoi(dbStr); err == nil {
			config.DB = db
		}
	}

	return config, nil
}
