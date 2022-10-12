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

type Cache interface {
	Connect(password string) error

	Get(key string) (string, error)
	Set(key string, value any, expiration time.Duration) (string, error)

	Properties() map[string]string
}

type cache struct {
	client  *redis.Client
	options redis.Options
}

func New(host string, port int) Cache {
	address := fmt.Sprintf("%s:%d", host, port)

	options := redis.Options{
		Addr: address,
		DB:   0,
	}

	return &cache{
		options: options,
	}
}

// Connect establishes a connection to the cache client.
// Returns an error if a connection cannot be established.
func (c *cache) Connect(password string) error {
	c.options.Password = password
	c.client = redis.NewClient(&c.options)

	_, err := c.client.Ping(ctx).Result()

	return err
}

// Returns a map of various cache connection properties.
func (c *cache) Properties() map[string]string {
	host, port, _ := net.SplitHostPort(c.options.Addr)

	return map[string]string{
		"host": host,
		"port": port,
	}
}

// Retrieves the value corresponding with the key.
// Returns ErrKeyNotFound if the key is not in cache.
func (c *cache) Get(key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
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
func (c *cache) Set(key string, value any, expiration time.Duration) (string, error) {
	return c.client.Set(ctx, key, value, expiration).Result()
}
