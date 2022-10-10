package cache

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var ErrKeyNotFound = errors.New("cache: no matching key found")

type Cache struct {
	client  *redis.Client
	options redis.Options
}

func New(host string, port int) Cache {
	address := fmt.Sprintf("%s:%d", host, port)

	options := redis.Options{
		Addr: address,
		DB:   0,
	}

	return Cache{
		options: options,
	}
}

// Connect establishes a connection to the cache client.
// Returns an error if a connection cannot be established.
func (cache *Cache) Connect(password string) error {
	cache.options.Password = password
	cache.client = redis.NewClient(&cache.options)

	_, err := cache.client.Ping(ctx).Result()

	return err
}

// Returns a map of various cache connection properties.
func (cache *Cache) Properties() map[string]string {
	host, port, _ := net.SplitHostPort(cache.options.Addr)

	return map[string]string{
		"host": host,
		"port": port,
	}
}

// Retrieves the value corresponding with the key.
// Returns ErrKeyNotFound if the key is not in cache.
func (cache *Cache) Get(key string) (string, error) {
	value, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrKeyNotFound
		}

		return "", err
	}

	return value, nil
}

// Sets the key value pair in cache with an expiration TTL.
// An expiration value of 0 indicates no expiration.
func (cache *Cache) Set(key string, value any, expiration time.Duration) (string, error) {
	return cache.client.Set(ctx, key, value, expiration).Result()
}
