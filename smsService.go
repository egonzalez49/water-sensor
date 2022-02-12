package main

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func sendSms() {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	apiKey := os.Getenv("TWILIO_API_KEY")
	apiSecret := os.Getenv("TWILIO_API_SECRET")

	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username:   apiKey,
		Password:   apiSecret,
		AccountSid: accountSid,
	})

	params := &openapi.CreateMessageParams{}

	toPhoneNum := os.Getenv("TO_PHONE_NUM")
	fromPhoneNum := os.Getenv("FROM_PHONE_NUM")
	params.SetTo(toPhoneNum)
	params.SetFrom(fromPhoneNum)
	params.SetBody("Water leak detected.")

	resp, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		err = nil
	} else {
		fmt.Println("Message Sid: " + *resp.Sid)
	}
}
