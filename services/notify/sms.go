package notify

import (
	"encoding/json"

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

	// Send an sms to each recipient in concurrent-fashion.
	for _, num := range recipientNumbers {
		go func(num string) {
			n._sendSms(client, senderNumber, num)
		}(num)
	}
}

func (n *Notifier) _sendSms(client *twilio.RestClient, senderNumber string, recipientNumber string) {
	params := &openapi.CreateMessageParams{}

	params.SetTo(senderNumber)
	params.SetFrom(recipientNumber)
	params.SetBody("Water leak detected.")

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		n.Logger.Errorf("failed to send SMS to %s\n", recipientNumber)
		n.Logger.Errorf("error: %v\n", err)
	} else {
		n.Logger.Infof("successfully sent SMS to %s\n", recipientNumber)
		response, _ := json.Marshal(*resp)
		n.Logger.Infof("response: %s\n", response)
	}
}
