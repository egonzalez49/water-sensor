package mocks

import (
	"time"

	"github.com/egonzalez49/water-sensor/internal/cache"
)

const ExistingCacheKey = "999"

type MockCache struct {
	Keys []string
}

func (c *MockCache) Connect(password string) error {
	return nil
}

func (c *MockCache) Get(key string) (string, error) {
	if key == ExistingCacheKey {
		return "1", nil
	}

	return "", cache.ErrKeyNotFound
}

func (c *MockCache) Set(key string, value any, expiration time.Duration) (string, error) {
	c.Keys = append(c.Keys, key)

	return "OK", nil
}

func (c *MockCache) Properties() map[string]string {
	return nil
}
