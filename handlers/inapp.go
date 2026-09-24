package handlers

import (
	"fmt"

	"notification-agent/models"
)

type InAppHandler struct{}

func (InAppHandler) Send(notification models.Notification) error {
	fmt.Printf(
		"[IN-APP] notification=%s user=%s message=%s\n",
		notification.ID,
		notification.UserID,
		notification.Message,
	)

	return nil
}
