package config

import "os"

type MqttConfig struct {
	ClientId string
	Host     string
	Port     string
	Username string
	Password string
}

func loadMqttConfig() MqttConfig {
	return MqttConfig{
		ClientId: os.Getenv("MQTT_CLID"),
		Host:     os.Getenv("MQTT_HOST"),
		Port:     os.Getenv("MQTT_PORT"),
		Username: os.Getenv("MQTT_USER"),
		Password: os.Getenv("MQTT_PASS"),
	}
}
