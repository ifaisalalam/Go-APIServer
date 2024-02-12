package cache

import "fmt"

const (
	// Redis constant for getting Redis cache client.
	Redis = "redis"
)

// NewCache initializes the cache instance based on Config.
func NewCache(config *Config) (Cache, error) {
	switch config.Driver {
	case Redis:
		r, err := NewRedisClient(&config.RedisConfig)
		if err != nil {
			return nil, err
		}
		return r, nil
	default:
		return nil, fmt.Errorf("unknown cache driver: %s", config.Driver)
	}
}
