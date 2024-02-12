package cache

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/go-redis/redis/v7"
)

// RedisConfig holds all required details for initializing Redis driver.
type RedisConfig struct {
	Host     string
	Port     string
	Database int32
	Password string

	hooks []redis.Hook
}

func (c *RedisConfig) WithHooks(hooks ...redis.Hook) *RedisConfig {
	for _, hook := range hooks {
		c.hooks = append(c.hooks, hook)
	}
	return c
}

// NewRedisClient initializes and returns a RedisClient instance.
func NewRedisClient(config *RedisConfig) (*RedisClient, error) {
	options := &redis.Options{
		Addr:     net.JoinHostPort(config.Host, config.Port),
		Password: config.Password,
		DB:       int(config.Database),
	}

	client := redis.NewClient(options)
	for _, hook := range config.hooks {
		client.AddHook(hook)
	}

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{client}, nil
}

// RedisClient implements the Cache interface for Redis.
type RedisClient struct {
	client *redis.Client
}

// Get value of a key from cache.
func (r *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNil
	}
	if err != nil {
		return nil, err
	}

	return []byte(val), nil
}

// Set a key in cache with TTL.
func (r *RedisClient) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := r.client.Set(key, string(value), ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

// Delete a key in cache.
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	err := r.client.Del(key).Err()
	return err
}

// IsAlive performs a healthcheck on the cache.
func (r *RedisClient) IsAlive(ctx context.Context) (bool, error) {
	ping, err := r.client.Ping().Result()
	if err != nil {
		return false, err
	}
	if ping == "PONG" {
		return true, err
	}
	return false, err
}
