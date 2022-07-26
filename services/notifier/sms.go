package notifier

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func executeSms() {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	apiKey := os.Getenv("TWILIO_API_KEY")
	apiSecret := os.Getenv("TWILIO_API_SECRET")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   apiKey,
		Password:   apiSecret,
		AccountSid: accountSid,
	})

	senderNumber := os.Getenv("SENDER_NUMBER")
	recipientNumbers := parseRecipientNumbers()

	for i := range recipientNumbers {
		sendSms(client, senderNumber, recipientNumbers[i])
	}
}

func sendSms(client *twilio.RestClient, senderNumber string, recipientNumber string) {
	params := &openapi.CreateMessageParams{}

	params.SetTo(senderNumber)
	params.SetFrom(recipientNumber)
	params.SetBody("Water leak detected.")

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		message := fmt.Sprintf("Failed to send SMS to %s.", recipientNumber)
		fmt.Println(message + " Error: " + err.Error())
	} else {
		message := fmt.Sprintf("Successfully sent SMS to %s.", recipientNumber)
		response, _ := json.Marshal(*resp)
		fmt.Println(message + " Response: " + string(response))
	}
}

func parseRecipientNumbers() []string {
	return strings.Split(strings.ReplaceAll(os.Getenv("RECIPIENT_NUMBERS"), " ", ""), ",")
}
