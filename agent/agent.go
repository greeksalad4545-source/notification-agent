package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"notification-agent/models"
	"notification-agent/processor"
	"notification-agent/servicebus"
	"notification-agent/storage"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type MessageReceiver interface {
	Receive(ctx context.Context) (*azservicebus.ReceivedMessage, error)
	Complete(ctx context.Context, message *azservicebus.ReceivedMessage) error
	Abandon(ctx context.Context, message *azservicebus.ReceivedMessage) error
}

type NotificationStore interface {
	SaveNotification(
		ctx context.Context,
		notification models.Notification,
		status string,
	) error
}

type Agent struct {
	consumer  MessageReceiver
	db        NotificationStore
	processor *processor.Processor
}

func New(
	consumer MessageReceiver,
	db NotificationStore,
	notificationProcessor *processor.Processor,
) *Agent {
	return &Agent{
		consumer:  consumer,
		db:        db,
		processor: notificationProcessor,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			fmt.Println("Agent shutdown requested.")
			return nil
		}

		err := a.ProcessNext(ctx)
		if err != nil {
			fmt.Println("Agent processing error:", err)
		}
	}
}

func (a *Agent) ProcessNext(ctx context.Context) error {
	fmt.Println("Waiting for notification...")

	message, err := a.consumer.Receive(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}

		return fmt.Errorf("failed to receive message: %w", err)
	}

	if message == nil {
		return nil
	}

	fmt.Println("Message received from Azure Service Bus:")
	fmt.Println(string(message.Body))

	var notification models.Notification

	err = json.Unmarshal(message.Body, &notification)
	if err != nil {
		fmt.Println("Failed to parse notification:", err)

		if abandonErr := a.consumer.Abandon(ctx, message); abandonErr != nil {
			return fmt.Errorf(
				"failed to abandon invalid message: %w",
				abandonErr,
			)
		}

		fmt.Println("Message abandoned and will be retried.")
		return fmt.Errorf("invalid notification JSON: %w", err)
	}

	fmt.Println("Notification parsed successfully:")
	fmt.Printf("%+v\n", notification)

	err = a.processor.Process(notification)
	if err != nil {
		fmt.Printf(
			"Error processing notification %s: %v\n",
			notification.ID,
			err,
		)

		if abandonErr := a.consumer.Abandon(ctx, message); abandonErr != nil {
			return fmt.Errorf(
				"failed to abandon processing-failed message: %w",
				abandonErr,
			)
		}

		fmt.Println("Message abandoned and will be retried.")
		return fmt.Errorf("failed to process notification: %w", err)
	}

	fmt.Printf(
		"Notification %s processed successfully.\n",
		notification.ID,
	)

	err = a.db.SaveNotification(
		ctx,
		notification,
		"processed",
	)
	if err != nil {
		fmt.Println("Failed to save notification:", err)

		if abandonErr := a.consumer.Abandon(ctx, message); abandonErr != nil {
			return fmt.Errorf(
				"failed to abandon database-failed message: %w",
				abandonErr,
			)
		}

		fmt.Println("Message abandoned and will be retried.")
		return fmt.Errorf("failed to save notification: %w", err)
	}

	fmt.Println("Notification saved to PostgreSQL.")

	err = a.consumer.Complete(ctx, message)
	if err != nil {
		return fmt.Errorf(
			"failed to complete message: %w",
			err,
		)
	}

	fmt.Println("Message completed successfully.")

	return nil
}

var _ MessageReceiver = (*servicebus.Consumer)(nil)
var _ NotificationStore = (*storage.Postgres)(nil)
