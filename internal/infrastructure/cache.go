package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/config"

	"github.com/redis/go-redis/v9"
)

// CacheInterface defines the interface for cache operations
type CacheInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	CheckRateLimit(ctx context.Context, identifier string, maxRequests int64, window time.Duration) (bool, error)
	Close() error
}

type Cache struct {
	client *redis.Client
}

func NewCache(cfg *config.Config) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPwd,
		DB:       cfg.RedisDB,
	})
	return &Cache{client: client}
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

func (c *Cache) CheckRateLimit(ctx context.Context, identifier string, maxRequests int64, window time.Duration) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", identifier)
	count, err := c.Increment(ctx, key)
	if err != nil {
		return false, err
	}
	if count == 1 {
		if err := c.Expire(ctx, key, window); err != nil {
			return false, err
		}
	}
	return count <= maxRequests, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}
