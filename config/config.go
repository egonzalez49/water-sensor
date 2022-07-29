package config

import (
	"errors"

	"github.com/joho/godotenv"
)

type Config struct {
	Log    LogConfig
	Mqtt   MqttConfig
	Redis  RedisConfig
	Twilio TwilioConfig
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("error loading .env file")
	}

	cfg := &Config{
		Log:    loadLogConfig(),
		Mqtt:   loadMqttConfig(),
		Redis:  loadRedisConfig(),
		Twilio: loadTwilioConfig(),
	}

	return cfg, nil
}
