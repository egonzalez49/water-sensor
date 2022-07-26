package config

import (
	"errors"

	"github.com/joho/godotenv"
)

type Config struct {
	Mqtt  MqttConfig
	Redis RedisConfig
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("error loading .env file")
	}

	cfg := &Config{
		Mqtt:  loadMqttConfig(),
		Redis: loadRedisConfig(),
	}

	return cfg, nil
}
