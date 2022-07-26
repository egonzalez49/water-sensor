package config

import (
	"os"
	"strings"
)

type TwilioConfig struct {
	Sid    string
	Key    string
	Secret string

	SenderNumber     string
	RecipientNumbers string
}

func loadTwilioConfig() TwilioConfig {
	return TwilioConfig{
		Sid:              os.Getenv("TWILIO_SID"),
		Key:              os.Getenv("TWILIO_KEY"),
		Secret:           os.Getenv("TWILIO_SECRET"),
		SenderNumber:     os.Getenv("SENDER_NUMBER"),
		RecipientNumbers: os.Getenv("RECIPIENT_NUMBERS"),
	}
}

func (t *TwilioConfig) GetRecipientNumbers() []string {
	return strings.Split(strings.ReplaceAll(os.Getenv("RECIPIENT_NUMBERS"), " ", ""), ",")
}
