package handlers

import (
	"fmt"

	"notification-agent/models"
)

type EmailHandler struct{}

func (EmailHandler) Send(notification models.Notification) error {
	fmt.Printf(
		"[EMAIL] notification=%s user=%s subject=%s message=%s\n",
		notification.ID,
		notification.UserID,
		notification.Subject,
		notification.Message,
	)

	return nil
}
