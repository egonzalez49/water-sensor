package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/egonzalez49/water-sensor/internal/broker"
	"github.com/egonzalez49/water-sensor/internal/cache"
	"github.com/egonzalez49/water-sensor/internal/logger"
	"github.com/egonzalez49/water-sensor/internal/sms"
)

const envPrefix = "SENSOR"

type application struct {
	config config
	logger *logger.Logger
	cache  cache.Cache
	broker broker.Broker
	sms    sms.Sms
}

func main() {
	var cfg config
	err := loadConfig("../../", &cfg)
	if err != nil {
		log.Fatalf("unable to load configuration: %s\n", err)
	}

	// Determine what to set the logger level to.
	// If no value is set from configuration,
	// default the log level to INFO.
	var logLevel logger.Level
	if cfg.LogLevel != "" {
		logLevel = logger.LevelFromString(cfg.LogLevel)
	} else {
		logLevel = logger.LevelInfo
	}

	logger := logger.New(os.Stdout, logLevel)

	app := application{
		config: cfg,
		logger: logger,
		cache:  cache.New(cfg.Redis.Host, cfg.Redis.Port),
		broker: broker.New(cfg.Mqtt.ClientId, cfg.Mqtt.Host, cfg.Mqtt.Port),
		sms:    sms.New(cfg.Twilio.Sid, cfg.Twilio.Key, cfg.Twilio.Secret, cfg.Twilio.Sender),
	}

	app.connectToCache()
	app.connectToBroker()

	// Keep the main goroutine running by
	// listening for a shutdown signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	app.broker.Disconnect(5000)
	app.logger.Info("shutting down service", map[string]string{
		"signal": s.String(),
	})
}

// Establishes a connection to the application's configured cache.
// If one cannot be established, logs the error and exits the program.
func (app *application) connectToCache() {
	err := app.cache.Connect(app.config.Redis.Password)
	if err != nil {
		app.logger.Fatal(err, nil)
	}

	app.logger.Info("connected to cache", app.cache.Properties())
}

// Establishes a connection to the application's configured broker.
// If one cannot be established, logs the error and exits the program.
func (app *application) connectToBroker() {
	// Set various broker message handlers
	app.broker.SetOnConnectHandler(app.onBrokerConnection)
	app.broker.SetOnConnectionLostHandler(app.onBrokerConnectionLost)
	app.broker.SetDefaultMessageHandler(app.onBrokerUnknownMessage)

	err := app.broker.Connect(app.config.Mqtt.Username, app.config.Mqtt.Password)
	if err != nil {
		app.logger.Fatal(err, nil)
	}

	// Each key represents a broker topic to subscribe to.
	// Each value represents the quality of service level.
	// 0 - at most once
	// 1 - at least once
	// 2 - exactly once
	filters := map[string]byte{"sensor/water": 1}
	app.broker.Subscribe(filters, app.onWaterSensorHandler)
}
