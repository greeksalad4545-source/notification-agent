package agent

import (
	"context"
	"errors"
	"testing"

	"notification-agent/handlers"
	"notification-agent/models"
	"notification-agent/processor"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type fakeMessageReceiver struct {
	message       *azservicebus.ReceivedMessage
	receiveErr    error
	completeErr   error
	abandonErr    error
	completeCalls int
	abandonCalls  int
}

func (f *fakeMessageReceiver) Receive(
	ctx context.Context,
) (*azservicebus.ReceivedMessage, error) {
	if f.receiveErr != nil {
		return nil, f.receiveErr
	}

	return f.message, nil
}

func (f *fakeMessageReceiver) Complete(
	ctx context.Context,
	message *azservicebus.ReceivedMessage,
) error {
	f.completeCalls++
	return f.completeErr
}

func (f *fakeMessageReceiver) Abandon(
	ctx context.Context,
	message *azservicebus.ReceivedMessage,
) error {
	f.abandonCalls++
	return f.abandonErr
}

type fakeNotificationStore struct {
	saveCalls    int
	savedStatus  string
	savedMessage models.Notification
	err          error
}

func (f *fakeNotificationStore) SaveNotification(
	ctx context.Context,
	notification models.Notification,
	status string,
) error {
	f.saveCalls++
	f.savedMessage = notification
	f.savedStatus = status
	return f.err
}

func newTestProcessor() *processor.Processor {
	return processor.New(
		map[string]handlers.NotificationHandler{
			"email":  handlers.EmailHandler{},
			"sms":    handlers.SMSHandler{},
			"in-app": handlers.InAppHandler{},
		},
	)
}

func TestProcessNextSuccess(t *testing.T) {
	receiver := &fakeMessageReceiver{
		message: &azservicebus.ReceivedMessage{
			Body: []byte(`{
				"id": "agent-test-001",
				"user_id": "user-123",
				"channel": "email",
				"subject": "Test notification",
				"message": "Hello from agent test"
			}`),
		},
	}

	store := &fakeNotificationStore{}

	agent := New(
		receiver,
		store,
		newTestProcessor(),
	)

	err := agent.ProcessNext(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if store.saveCalls != 1 {
		t.Fatalf("expected 1 save call, got %d", store.saveCalls)
	}

	if store.savedStatus != "processed" {
		t.Errorf(
			"expected status processed, got %s",
			store.savedStatus,
		)
	}

	if store.savedMessage.ID != "agent-test-001" {
		t.Errorf(
			"expected notification ID agent-test-001, got %s",
			store.savedMessage.ID,
		)
	}

	if receiver.completeCalls != 1 {
		t.Fatalf(
			"expected 1 complete call, got %d",
			receiver.completeCalls,
		)
	}

	if receiver.abandonCalls != 0 {
		t.Fatalf(
			"expected 0 abandon calls, got %d",
			receiver.abandonCalls,
		)
	}
}

func TestProcessNextInvalidJSON(t *testing.T) {
	receiver := &fakeMessageReceiver{
		message: &azservicebus.ReceivedMessage{
			Body: []byte(`not valid json`),
		},
	}

	store := &fakeNotificationStore{}

	agent := New(
		receiver,
		store,
		newTestProcessor(),
	)

	err := agent.ProcessNext(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if store.saveCalls != 0 {
		t.Fatalf(
			"expected 0 save calls, got %d",
			store.saveCalls,
		)
	}

	if receiver.completeCalls != 0 {
		t.Fatalf(
			"expected 0 complete calls, got %d",
			receiver.completeCalls,
		)
	}

	if receiver.abandonCalls != 1 {
		t.Fatalf(
			"expected 1 abandon call, got %d",
			receiver.abandonCalls,
		)
	}
}

func TestProcessNextUnsupportedChannel(t *testing.T) {
	receiver := &fakeMessageReceiver{
		message: &azservicebus.ReceivedMessage{
			Body: []byte(`{
				"id": "agent-test-unsupported",
				"user_id": "user-123",
				"channel": "telegram",
				"message": "Unsupported channel"
			}`),
		},
	}

	store := &fakeNotificationStore{}

	agent := New(
		receiver,
		store,
		newTestProcessor(),
	)

	err := agent.ProcessNext(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if store.saveCalls != 0 {
		t.Fatalf(
			"expected 0 save calls, got %d",
			store.saveCalls,
		)
	}

	if receiver.completeCalls != 0 {
		t.Fatalf(
			"expected 0 complete calls, got %d",
			receiver.completeCalls,
		)
	}

	if receiver.abandonCalls != 1 {
		t.Fatalf(
			"expected 1 abandon call, got %d",
			receiver.abandonCalls,
		)
	}
}

func TestProcessNextDatabaseError(t *testing.T) {
	receiver := &fakeMessageReceiver{
		message: &azservicebus.ReceivedMessage{
			Body: []byte(`{
				"id": "agent-test-db-error",
				"user_id": "user-123",
				"channel": "email",
				"subject": "Database test",
				"message": "Database failure"
			}`),
		},
	}

	store := &fakeNotificationStore{
		err: errors.New("database unavailable"),
	}

	agent := New(
		receiver,
		store,
		newTestProcessor(),
	)

	err := agent.ProcessNext(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if store.saveCalls != 1 {
		t.Fatalf(
			"expected 1 save call, got %d",
			store.saveCalls,
		)
	}

	if receiver.completeCalls != 0 {
		t.Fatalf(
			"expected 0 complete calls, got %d",
			receiver.completeCalls,
		)
	}

	if receiver.abandonCalls != 1 {
		t.Fatalf(
			"expected 1 abandon call, got %d",
			receiver.abandonCalls,
		)
	}
}

func TestProcessNextReceiveError(t *testing.T) {
	receiver := &fakeMessageReceiver{
		receiveErr: errors.New("service bus unavailable"),
	}

	store := &fakeNotificationStore{}

	agent := New(
		receiver,
		store,
		newTestProcessor(),
	)

	err := agent.ProcessNext(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if store.saveCalls != 0 {
		t.Fatalf(
			"expected 0 save calls, got %d",
			store.saveCalls,
		)
	}

	if receiver.completeCalls != 0 {
		t.Fatalf(
			"expected 0 complete calls, got %d",
			receiver.completeCalls,
		)
	}

	if receiver.abandonCalls != 0 {
		t.Fatalf(
			"expected 0 abandon calls, got %d",
			receiver.abandonCalls,
		)
	}
}

func TestProcessNextNoMessage(t *testing.T) {
	receiver := &fakeMessageReceiver{}

	store := &fakeNotificationStore{}

	agent := New(
		receiver,
		store,
		newTestProcessor(),
	)

	err := agent.ProcessNext(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if store.saveCalls != 0 {
		t.Fatalf(
			"expected 0 save calls, got %d",
			store.saveCalls,
		)
	}
}
