package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/egonzalez49/water-sensor/config"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	redis *redis.Client
}

func NewCache(cfg *config.Config) *Cache {
	return &Cache{
		redis: initRedis(cfg),
	}
}

func initRedis(cfg *config.Config) *redis.Client {
	host := cfg.Redis.Host
	port := cfg.Redis.Port
	address := fmt.Sprintf("%s:%s", host, port)

	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.Redis.Password,
		DB:       0,
	})
}

func (cache *Cache) Get(ctx context.Context, key string) (string, error) {
	return cache.redis.Get(ctx, key).Result()
}

func (cache *Cache) Set(
	ctx context.Context, key string, value interface{}, expiration time.Duration,
) (string, error) {
	return cache.redis.Set(ctx, key, value, expiration).Result()
}
