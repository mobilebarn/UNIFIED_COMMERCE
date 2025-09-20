package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"unified-commerce/services/shared/database"
	"unified-commerce/services/shared/logger"
)

// CacheService provides high-level caching operations using Redis
type CacheService struct {
	client *database.RedisClient
	logger *logger.Logger
	prefix string
}

// NewCacheService creates a new cache service
func NewCacheService(client *database.RedisClient, logger *logger.Logger, serviceName string) *CacheService {
	return &CacheService{
		client: client,
		logger: logger,
		prefix: fmt.Sprintf("%s:", serviceName),
	}
}

// Get retrieves a value from cache and unmarshals it into the provided interface
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	if c.client == nil {
		return fmt.Errorf("redis client not available")
	}

	fullKey := c.prefix + key
	data, err := c.client.Get(ctx, fullKey)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Set stores a value in cache with expiration
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client not available")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	fullKey := c.prefix + key
	return c.client.Set(ctx, fullKey, data, expiration)
}

// Delete removes a key from cache
func (c *CacheService) Delete(ctx context.Context, keys ...string) error {
	if c.client == nil {
		return fmt.Errorf("redis client not available")
	}

	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = c.prefix + key
	}

	return c.client.Delete(ctx, fullKeys...)
}

// Exists checks if a key exists in cache
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("redis client not available")
	}

	fullKey := c.prefix + key
	count, err := c.client.Exists(ctx, fullKey)
	return count > 0, err
}

// SetExpiration sets expiration for an existing key
func (c *CacheService) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client not available")
	}

	fullKey := c.prefix + key
	return c.client.Expire(ctx, fullKey, expiration)
}

// Increment increments a numeric value in cache
func (c *CacheService) Increment(ctx context.Context, key string) (int64, error) {
	if c.client == nil {
		return 0, fmt.Errorf("redis client not available")
	}

	fullKey := c.prefix + key
	return c.client.Increment(ctx, fullKey)
}

// IncrementBy increments a numeric value by a specific amount
func (c *CacheService) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	if c.client == nil {
		return 0, fmt.Errorf("redis client not available")
	}

	fullKey := c.prefix + key
	return c.client.IncrementBy(ctx, fullKey, value)
}

// GetOrSet retrieves a value from cache, or sets it using the provided function if not found
func (c *CacheService) GetOrSet(ctx context.Context, key string, dest interface{}, expiration time.Duration, fetcher func() (interface{}, error)) error {
	// Try to get from cache first
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil // Found in cache
	}

	// Not in cache, fetch the data
	value, err := fetcher()
	if err != nil {
		return err
	}

	// Store in cache
	if err := c.Set(ctx, key, value, expiration); err != nil {
		c.logger.WithError(err).Warn("Failed to store value in cache")
	}

	// Marshal the value into dest
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// SessionCache provides session-specific caching operations
type SessionCache struct {
	cache     *CacheService
	sessionID string
}

// NewSessionCache creates a session-specific cache
func (c *CacheService) NewSessionCache(sessionID string) *SessionCache {
	return &SessionCache{
		cache:     c,
		sessionID: sessionID,
	}
}

// Get retrieves a session-specific value
func (s *SessionCache) Get(ctx context.Context, key string, dest interface{}) error {
	sessionKey := fmt.Sprintf("session:%s:%s", s.sessionID, key)
	return s.cache.Get(ctx, sessionKey, dest)
}

// Set stores a session-specific value
func (s *SessionCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	sessionKey := fmt.Sprintf("session:%s:%s", s.sessionID, key)
	return s.cache.Set(ctx, sessionKey, value, expiration)
}

// Delete removes session-specific keys
func (s *SessionCache) Delete(ctx context.Context, keys ...string) error {
	sessionKeys := make([]string, len(keys))
	for i, key := range keys {
		sessionKeys[i] = fmt.Sprintf("session:%s:%s", s.sessionID, key)
	}
	return s.cache.Delete(ctx, sessionKeys...)
}

// ClearSession removes all session data
func (s *SessionCache) ClearSession(ctx context.Context) error {
	sessionPattern := fmt.Sprintf("session:%s:*", s.sessionID)
	// Note: This is a simplified implementation. In production, you might want to use SCAN
	// to avoid blocking operations
	return s.cache.Delete(ctx, sessionPattern)
}
