package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
)

func (app *application) messageProperties(msg mqtt.Message) map[string]string {
	payload := msg.Payload()
	return map[string]string{
		"topic":   msg.Topic(),
		"payload": string(payload),
	}
}

func (app *application) onBrokerConnection(client mqtt.Client) {
	app.logger.Info("connection to broker established", app.broker.Properties())
}

func (app *application) onBrokerConnectionLost(client mqtt.Client, err error) {
	app.logger.Error(fmt.Errorf("connection to broker lost: %w", err), app.broker.Properties())
}

func (app *application) onBrokerUnknownMessage(client mqtt.Client, msg mqtt.Message) {
	app.logger.Info("unknown message received", app.messageProperties(msg))
}

type MessagePayload struct {
	Time  string
	Model string
	Id    string
	Event string
	Code  string
	Mic   string
}

func (app *application) onWaterSensorHandler(client mqtt.Client, msg mqtt.Message) {
	app.logger.Info("received message from broker", app.messageProperties(msg))

	var data MessagePayload
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		app.logger.Error(fmt.Errorf("unable to unmarshal message: %w", err), app.messageProperties(msg))
		return
	}

	ctx := context.Background()
	_, err := app.cache.Get(ctx, data.Id)
	if err == redis.Nil {
		// Key does not exist in cache.
		// Notify respective parties and save the id
		// in cache to prevent processing duplicates.
		app.sms.Send(app.config.Twilio.Recipients, "Water leak detected.")

		_, err = app.cache.Set(ctx, data.Id, struct{}{}, 5*time.Minute)
		if err != nil {
			app.logger.Error(fmt.Errorf("unexpected error with cache: %w", err), nil)
			return
		}
	} else if err != nil {
		app.logger.Error(fmt.Errorf("unexpected error with cache: %w", err), nil)
		return
	}
}
