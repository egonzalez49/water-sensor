package config

import "os"

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
	}
}
