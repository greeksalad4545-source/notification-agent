package processor

import (
	"fmt"

	"notification-agent/handlers"
	"notification-agent/models"
)

func Process(notification models.Notification) error {

	var handler handlers.NotificationHandler

	switch notification.Channel {

	case "email":
		handler = handlers.EmailHandler{}

	case "sms":
		handler = handlers.SMSHandler{}

	case "in-app":
		handler = handlers.InAppHandler{}

	default:
		return fmt.Errorf("unsupported notification channel: %s", notification.Channel)
	}

	return handler.Send(notification)
}
