package sms

import (
	"encoding/json"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type MessageSender interface {
	SendMessage(recipientNumber string, body string) error
}

type Sms struct {
	client *twilio.RestClient
	sender string
}

func New(sid, key, secret, sender string) Sms {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   key,
		Password:   secret,
		AccountSid: sid,
	})

	return Sms{
		client: client,
		sender: sender,
	}
}

func (s *Sms) Send(recipients []string, body string) {
	// Send an sms to each recipient in concurrent-fashion.
	for _, num := range recipients {
		go func(num string) {
			s.sendMessage(num, body)
		}(num)
	}
}

func (s *Sms) sendMessage(recipient, body string) ([]byte, error) {
	params := &openapi.CreateMessageParams{}
	params.SetTo(s.sender)
	params.SetFrom(recipient)
	params.SetBody(body)

	response, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return nil, err
	}

	jsonResponse, err := json.Marshal(response)
	return jsonResponse, err
}
