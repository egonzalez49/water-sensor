package subscriber

import (
	"context"
	"encoding/json"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/egonzalez49/water-sensor/config"
	"github.com/egonzalez49/water-sensor/services/cache"
	"github.com/egonzalez49/water-sensor/services/notifier"
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
	Config *config.Config
	Cache  *cache.Cache
}

func NewSubscriber(cfg *config.Config, inmem *cache.Cache) *Service {
	return &Service{
		Config: cfg,
		Cache:  inmem,
	}
}

func (s *Service) OnConnect(client mqtt.Client) {
	log.Println("Connected to broker.")
}

func (s *Service) OnConnectionLost(client mqtt.Client, err error) {
	log.Printf("Connection to broker lost: %v.\n", err)
}

func (s *Service) OnUnknownMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Unknown message from topic %s received: %s\n", msg.Topic(), msg.Payload())
}

func (s *Service) OnWaterSensorHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Message from topic %s received: %s.\n", msg.Topic(), msg.Payload())

	var data SensorPayload
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Unable to unmarshal message: %v\n", err)
		return
	}

	_, err := s.Cache.Get(context.Background(), data.Id)
	if err == redis.Nil {
		// Key does not exist in cache.
		s.Cache.Set(context.Background(), data.Id, struct{}{}, 5*time.Minute)
		notifier.Notify()
	} else if err != nil {
		log.Printf("Redis error: %v\n", err)
		return
	}
}
