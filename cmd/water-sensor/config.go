package main

import (
	"github.com/spf13/viper"
)

type config struct {
	LogLevel string `mapstructure:"SENSOR_LOG_LEVEL"`

	Mqtt   mqttConfig   `mapstructure:",squash"`
	Redis  redisConfig  `mapstructure:",squash"`
	Twilio twilioConfig `mapstructure:",squash"`
}

type mqttConfig struct {
	ClientId string `mapstructure:"SENSOR_MQTT_CLID"`
	Host     string `mapstructure:"SENSOR_MQTT_HOST"`
	Port     int    `mapstructure:"SENSOR_MQTT_PORT"`
	Username string `mapstructure:"SENSOR_MQTT_USER"`
	Password string `mapstructure:"SENSOR_MQTT_PASS"`
}

type redisConfig struct {
	Host     string `mapstructure:"SENSOR_REDIS_HOST"`
	Port     int    `mapstructure:"SENSOR_REDIS_PORT"`
	Password string `mapstructure:"SENSOR_REDIS_PASS"`
}

type twilioConfig struct {
	Sid        string   `mapstructure:"SENSOR_TWILIO_SID"`
	Key        string   `mapstructure:"SENSOR_TWILIO_KEY"`
	Secret     string   `mapstructure:"SENSOR_TWILIO_SECRET"`
	Sender     string   `mapstructure:"SENSOR_SENDER_NUMBER"`
	Recipients []string `mapstructure:"SENSOR_RECIPIENT_NUMBERS"`
}

func loadConfig(path string, dest *config) error {
	// Tell Viper where to look for the config file.
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	// Expect all environment variables to
	// contain the given prefix.
	viper.SetEnvPrefix(envPrefix)

	// Automatically override any config values
	// if an equivalent environment variable exists.
	viper.AutomaticEnv()

	// Read the values from the config.
	// Ignore any errors due to config file not being found.
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Unmarshal any config or environment values
	// into the given config struct based on struct tags.
	err = viper.Unmarshal(dest)
	if err != nil {
		return err
	}

	return nil
}
