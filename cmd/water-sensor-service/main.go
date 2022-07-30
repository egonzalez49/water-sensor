package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/egonzalez49/water-sensor/config"
	"github.com/egonzalez49/water-sensor/logging"
	"github.com/egonzalez49/water-sensor/services/broker"
	"github.com/egonzalez49/water-sensor/services/cache"
	"github.com/egonzalez49/water-sensor/services/notify"
	"github.com/egonzalez49/water-sensor/services/subscriber"
)

func main() {
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

	filters := map[string]byte{"sensor/water": 1}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("can't initialize config: %v", err)
	}

	logger, err := logging.NewLogger(cfg)
	if err != nil {
		log.Fatalf("can't initialize logger: %v", err)
	}

	// Flush any remaining log buffer
	defer logger.Shutdown()

	inmem := cache.NewCache(cfg)
	notifier := notify.NewNotifier(cfg, logger)
	sub := subscriber.NewSubscriber(cfg, inmem, logger, notifier)

	bkr, err := broker.NewBroker(cfg, logger)
	if err != nil {
		logger.Fatalf("can't initialize broker: %v", err)
	}

	bkr.SetConnectionHandlers(sub.OnConnect, sub.OnConnectionLost)
	bkr.SetDefaultPublishHandler(sub.OnUnknownMessageHandler)

	err = bkr.Connect()
	if err != nil {
		logger.Fatalf("can't connect to broker: %v", err)
	}

	bkr.Subscribe(filters, sub.OnWaterSensorHandler)

	<-keepAlive
}
