package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"notification-agent/models"
	"notification-agent/processor"
	"notification-agent/servicebus"
	"notification-agent/storage"
)

type Agent struct {
	consumer *servicebus.Consumer
	db       *storage.Postgres
}

func New(
	consumer *servicebus.Consumer,
	db *storage.Postgres,
) *Agent {
	return &Agent{
		consumer: consumer,
		db:       db,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			fmt.Println("Agent shutdown requested.")
			return nil
		}

		fmt.Println("Waiting for notification...")

		message, err := a.consumer.Receive(ctx)
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("Stopping notification agent...")
				return nil
			}

			fmt.Println("Failed to receive message:", err)
			continue
		}

		if message == nil {
			continue
		}

		fmt.Println("Message received from Azure Service Bus:")
		fmt.Println(string(message.Body))

		var notification models.Notification

		err = json.Unmarshal(message.Body, &notification)
		if err != nil {
			fmt.Println("Failed to parse notification:", err)

			err = a.consumer.Abandon(ctx, message)
			if err != nil {
				fmt.Println("Failed to abandon message:", err)
			} else {
				fmt.Println("Message abandoned and will be retried.")
			}

			continue
		}

		fmt.Println("Notification parsed successfully:")
		fmt.Printf("%+v\n", notification)

		err = processor.Process(notification)
		if err != nil {
			fmt.Printf(
				"Error processing notification %s: %v\n",
				notification.ID,
				err,
			)

			err = a.consumer.Abandon(ctx, message)
			if err != nil {
				fmt.Println("Failed to abandon message:", err)
			} else {
				fmt.Println("Message abandoned and will be retried.")
			}

			continue
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

			err = a.consumer.Abandon(ctx, message)
			if err != nil {
				fmt.Println("Failed to abandon message:", err)
			} else {
				fmt.Println("Message abandoned and will be retried.")
			}

			continue
		}

		fmt.Println("Notification saved to PostgreSQL.")

		err = a.consumer.Complete(ctx, message)
		if err != nil {
			fmt.Println("Failed to complete message:", err)
			continue
		}

		fmt.Println("Message completed successfully.")
	}
}
