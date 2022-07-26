package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/egonzalez49/water-sensor/config"
	"github.com/egonzalez49/water-sensor/services/broker"
	"github.com/egonzalez49/water-sensor/services/cache"
	"github.com/egonzalez49/water-sensor/services/notify"
	"github.com/egonzalez49/water-sensor/services/subscriber"
)

func main() {
	topics := []string{"sensor/water"}

	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	inmem := cache.NewCache(cfg)
	notifier := notify.NewNotifier(cfg)
	sub := subscriber.NewSubscriber(cfg, inmem, notifier)

	bkr, err := broker.NewBroker(cfg)
	if err != nil {
		log.Fatal(err)
	}

	bkr.SetConnectionHandlers(sub.OnConnect, sub.OnConnectionLost)
	bkr.SetDefaultPublishHandler(sub.OnUnknownMessage)

	err = bkr.Connect()
	if err != nil {
		log.Fatal(err)
	}

	bkr.Subscribe(topics)

	<-keepAlive
}
