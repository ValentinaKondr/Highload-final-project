package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache wraps Redis client for caching metrics
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// Set stores a value with expiration
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// Get retrieves a value
func (r *RedisCache) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get value: %w", err)
	}
	return json.Unmarshal([]byte(val), dest)
}

// Increment increments a counter
func (r *RedisCache) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// GetInt64 retrieves an integer value
func (r *RedisCache) GetInt64(key string) (int64, error) {
	val, err := r.client.Get(r.ctx, key).Int64()
	if err == redis.Nil {
		return 0, fmt.Errorf("key not found")
	}
	return val, err
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}

