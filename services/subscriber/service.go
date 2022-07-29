package subscriber

import (
	"context"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/egonzalez49/water-sensor/config"
	"github.com/egonzalez49/water-sensor/logging"
	"github.com/egonzalez49/water-sensor/services/cache"
	"github.com/egonzalez49/water-sensor/services/notify"
	"github.com/go-redis/redis/v8"
)

type SensorPayload struct {
	Time  string
	Model string
	Id    string
	Event string
	Code  string
	Mic   string
}

type Service struct {
	Config   *config.Config
	Cache    *cache.Cache
	Logger   *logging.Logger
	Notifier *notify.Notifier
}

func NewSubscriber(cfg *config.Config, inmem *cache.Cache, logger *logging.Logger, notifier *notify.Notifier) *Service {
	return &Service{
		Config:   cfg,
		Cache:    inmem,
		Logger:   logger,
		Notifier: notifier,
	}
}

func (s *Service) OnConnect(client mqtt.Client) {
	s.Logger.Info("Connected to broker")
}

func (s *Service) OnConnectionLost(client mqtt.Client, err error) {
	s.Logger.Infof("Connection to broker lost: %v\n", err)
}

func (s *Service) OnUnknownMessageHandler(client mqtt.Client, msg mqtt.Message) {
	s.Logger.Infof("Unknown message from topic %s received: %s\n", msg.Topic(), msg.Payload())
}

func (s *Service) OnWaterSensorHandler(client mqtt.Client, msg mqtt.Message) {
	s.Logger.Infof("Message from topic %s received: %s\n", msg.Topic(), msg.Payload())

	var data SensorPayload
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		s.Logger.Errorf("Unable to unmarshal message: %v\n", err)
		return
	}

	ctx := context.Background()
	_, err := s.Cache.Get(ctx, data.Id)
	if err == redis.Nil {
		// Key does not exist in cache.
		// Notify respective parties and save the id
		// in cache to prevent processing duplicates.
		s.Notifier.Notify()

		_, err = s.Cache.Set(ctx, data.Id, struct{}{}, 5*time.Minute)
		if err != nil {
			s.Logger.Errorf("Redis error: %v\n", err)
			return
		}
	} else if err != nil {
		s.Logger.Errorf("Redis error: %v\n", err)
		return
	}
}
