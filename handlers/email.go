package handlers

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"notification-agent/models"
	"notification-agent/secrets"
)

type EmailHandler struct {
	keyVault    *secrets.KeyVault
	senderEmail string
}

func NewEmailHandler(
	keyVault *secrets.KeyVault,
	senderEmail string,
) *EmailHandler {
	return &EmailHandler{
		keyVault:    keyVault,
		senderEmail: senderEmail,
	}
}

func (h *EmailHandler) Send(
	notification models.Notification,
) error {
	if notification.RecipientEmail == "" {
		return fmt.Errorf("recipient email is required")
	}

	apiKey, err := h.keyVault.GetSecret(
		context.Background(),
		"SENDGRID-API-KEY",
	)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve SendGrid API key: %w",
			err,
		)
	}

	from := mail.NewEmail(
		"Notification Agent",
		h.senderEmail,
	)

	to := mail.NewEmail(
		notification.UserID,
		notification.RecipientEmail,
	)

	message := mail.NewSingleEmail(
		from,
		notification.Subject,
		to,
		notification.Message,
		notification.Message,
	)

	client := sendgrid.NewSendClient(apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf(
			"failed to send email: %w",
			err,
		)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf(
			"SendGrid returned status %d: %s",
			response.StatusCode,
			response.Body,
		)
	}

	fmt.Printf(
		"[EMAIL] notification=%s recipient=%s subject=%s status=%d\n",
		notification.ID,
		notification.RecipientEmail,
		notification.Subject,
		response.StatusCode,
	)

	return nil
}
