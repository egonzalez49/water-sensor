package sms

import (
	"encoding/json"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Sms interface {
	Send(recipient string, body string) ([]byte, error)
}

type sms struct {
	client *twilio.RestClient
	sender string
}

func New(sid, key, secret, sender string) Sms {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   key,
		Password:   secret,
		AccountSid: sid,
	})

	return &sms{
		client: client,
		sender: sender,
	}
}

func (s *sms) Send(recipient, body string) ([]byte, error) {
	params := &openapi.CreateMessageParams{}
	params.SetTo(s.sender).SetFrom(recipient).SetBody(body)

	response, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return nil, err
	}

	jsonResponse, err := json.Marshal(response)
	return jsonResponse, err
}

func (s *sms) sendMessage(recipient, body string) ([]byte, error) {
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
