package notify

import (
	"encoding/json"
	"log"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func (n *Notifier) sendSms() {
	accountSid := n.Config.Twilio.Sid
	apiKey := n.Config.Twilio.Key
	apiSecret := n.Config.Twilio.Secret

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   apiKey,
		Password:   apiSecret,
		AccountSid: accountSid,
	})

	senderNumber := n.Config.Twilio.SenderNumber
	recipientNumbers := n.Config.Twilio.GetRecipientNumbers()

	for i := range recipientNumbers {
		n._sendSms(client, senderNumber, recipientNumbers[i])
	}
}

func (n *Notifier) _sendSms(client *twilio.RestClient, senderNumber string, recipientNumber string) {
	params := &openapi.CreateMessageParams{}

	params.SetTo(senderNumber)
	params.SetFrom(recipientNumber)
	params.SetBody("Water leak detected.")

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS to %s\n", recipientNumber)
		log.Printf("Error: %v\n", err)
	} else {
		log.Printf("Successfully sent SMS to %s\n", recipientNumber)
		response, _ := json.Marshal(*resp)
		log.Printf("Response: %s\n", response)
	}
}
