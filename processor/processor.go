package processor

import (
	"fmt"

	"notification-agent/handlers"
	"notification-agent/models"
)

type Processor struct {
	handlers map[string]handlers.NotificationHandler
}

func New(
	handlers map[string]handlers.NotificationHandler,
) *Processor {
	return &Processor{
		handlers: handlers,
	}
}

func (p *Processor) Process(
	notification models.Notification,
) error {
	handler, ok := p.handlers[notification.Channel]
	if !ok {
		return fmt.Errorf(
			"unsupported notification channel: %s",
			notification.Channel,
		)
	}

	return handler.Send(notification)
}
