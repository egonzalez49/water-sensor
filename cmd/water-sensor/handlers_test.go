package main

import (
	"testing"

	"github.com/egonzalez49/water-sensor/internal/assert"
	"github.com/egonzalez49/water-sensor/internal/mocks"
)

func TestSendAlerts(t *testing.T) {
	tests := []struct {
		name         string
		message      messagePayload
		wantSmsCount int
	}{
		{
			name: "Non-existent Cache Key",
			message: messagePayload{
				Id: "123",
			},
			wantSmsCount: 2,
		},
		{
			name: "Existing Cache Key",
			message: messagePayload{
				Id: mocks.ExistingCacheKey,
			},
			wantSmsCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipients := []string{
				"+11234567890",
				"+11234567890",
			}
			cfg := newTestConfig(recipients)
			cache := &mocks.MockCache{}
			sms := &mocks.MockSms{}
			app := newTestApplication(t, cfg, cache, sms)

			// Because SMS alerts are sent through goroutines,
			// we want to wait for the calls to finish before verifying.
			sms.AddExpectedCalls(tt.wantSmsCount)
			app.sendAlerts(tt.message)
			sms.WaitForCalls()

			assert.Equal(t, sms.SendCount, tt.wantSmsCount)
		})
	}
}
