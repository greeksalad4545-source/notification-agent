package handlers

import (
	"context"
	"fmt"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"

	"notification-agent/models"
	"notification-agent/secrets"
)

type SMSHandler struct {
	keyVault *secrets.KeyVault
}

func NewSMSHandler(
	keyVault *secrets.KeyVault,
) *SMSHandler {
	return &SMSHandler{
		keyVault: keyVault,
	}
}

func (h *SMSHandler) Send(
	notification models.Notification,
) error {
	if notification.RecipientPhone == "" {
		return fmt.Errorf("recipient phone is required")
	}

	accountSID, err := h.keyVault.GetSecret(
		context.Background(),
		"TWILIO-ACCOUNT-SID",
	)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve Twilio account SID: %w",
			err,
		)
	}

	authToken, err := h.keyVault.GetSecret(
		context.Background(),
		"TWILIO-AUTH-TOKEN",
	)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve Twilio auth token: %w",
			err,
		)
	}

	fromNumber, err := h.keyVault.GetSecret(
		context.Background(),
		"TWILIO-FROM-NUMBER",
	)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve Twilio from number: %w",
			err,
		)
	}

	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: accountSID,
			Password: authToken,
		},
	)

	params := &api.CreateMessageParams{}
	params.SetTo(notification.RecipientPhone)
	params.SetFrom(fromNumber)
	params.SetBody(notification.Message)

	response, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf(
			"failed to send SMS: %w",
			err,
		)
	}

	fmt.Printf(
		"[SMS] notification=%s recipient=%s status=%s\n",
		notification.ID,
		notification.RecipientPhone,
		*response.Status,
	)

	return nil
}
