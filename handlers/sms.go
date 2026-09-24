package handlers

import (
	"fmt"

	"notification-agent/models"
)

type SMSHandler struct{}

func (SMSHandler) Send(notification models.Notification) error {
	fmt.Printf(
		"[SMS] notification=%s user=%s message=%s\n",
		notification.ID,
		notification.UserID,
		notification.Message,
	)

	return nil
}
