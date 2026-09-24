package servicebus

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Consumer struct {
	client   *azservicebus.Client
	receiver *azservicebus.Receiver
}

func NewConsumer(
	connectionString string,
	queueName string,
) (*Consumer, error) {
	client, err := azservicebus.NewClientFromConnectionString(
		connectionString,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Service Bus client: %w",
			err,
		)
	}

	receiver, err := client.NewReceiverForQueue(
		queueName,
		nil,
	)
	if err != nil {
		client.Close(context.Background())

		return nil, fmt.Errorf(
			"failed to create receiver: %w",
			err,
		)
	}

	return &Consumer{
		client:   client,
		receiver: receiver,
	}, nil
}

func (c *Consumer) Close(ctx context.Context) error {
	if err := c.receiver.Close(ctx); err != nil {
		return err
	}

	return c.client.Close(ctx)
}

func (c *Consumer) Receive(
	ctx context.Context,
) (*azservicebus.ReceivedMessage, error) {
	messages, err := c.receiver.ReceiveMessages(
		ctx,
		1,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to receive message: %w",
			err,
		)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

func (c *Consumer) Complete(
	ctx context.Context,
	message *azservicebus.ReceivedMessage,
) error {
	err := c.receiver.CompleteMessage(
		ctx,
		message,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to complete message: %w",
			err,
		)
	}

	return nil
}

func (c *Consumer) Abandon(
	ctx context.Context,
	message *azservicebus.ReceivedMessage,
) error {
	err := c.receiver.AbandonMessage(
		ctx,
		message,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to abandon message: %w",
			err,
		)
	}

	return nil
}
