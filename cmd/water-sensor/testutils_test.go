package main

import (
	"io"
	"testing"

	"github.com/egonzalez49/water-sensor/internal/cache"
	"github.com/egonzalez49/water-sensor/internal/logger"
	"github.com/egonzalez49/water-sensor/internal/sms"
)

func newTestConfig(recipients []string) config {
	return config{
		Twilio: twilioConfig{
			Recipients: recipients,
		},
	}
}

func newTestApplication(t *testing.T, cfg config, cache cache.Cache, sms sms.Sms) *application {
	return &application{
		config: cfg,
		cache:  cache,
		sms:    sms,
		logger: logger.New(io.Discard, logger.LevelDebug),
	}
}
