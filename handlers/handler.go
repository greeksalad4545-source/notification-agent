package handlers

import "notification-agent/models"

type NotificationHandler interface {
	Send(notification models.Notification) error
}
