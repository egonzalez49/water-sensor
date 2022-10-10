package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}

func New(host string, port int, password string) Cache {
	address := fmt.Sprintf("%s:%d", host, port)

	redis := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	return Cache{
		client: redis,
	}
}

func (cache *Cache) Get(ctx context.Context, key string) (string, error) {
	return cache.client.Get(ctx, key).Result()
}

func (cache *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) (string, error) {
	return cache.client.Set(ctx, key, value, expiration).Result()
}
