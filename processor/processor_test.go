package processor

import (
	"errors"
	"testing"

	"notification-agent/handlers"
	"notification-agent/models"
)

type fakeHandler struct {
	err   error
	calls int
}

func (f *fakeHandler) Send(notification models.Notification) error {
	f.calls++
	return f.err
}

func TestProcessEmail(t *testing.T) {
	emailHandler := &fakeHandler{}

	notificationProcessor := New(
		map[string]handlers.NotificationHandler{
			"email": emailHandler,
		},
	)

	notification := models.Notification{
		ID:      "test-email-001",
		UserID:  "user-123",
		Channel: "email",
		Subject: "Test Email",
		Message: "Hello from test",
	}

	err := notificationProcessor.Process(notification)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if emailHandler.calls != 1 {
		t.Fatalf(
			"expected handler to be called once, got %d",
			emailHandler.calls,
		)
	}
}

func TestProcessSMS(t *testing.T) {
	smsHandler := &fakeHandler{}

	notificationProcessor := New(
		map[string]handlers.NotificationHandler{
			"sms": smsHandler,
		},
	)

	notification := models.Notification{
		ID:      "test-sms-001",
		UserID:  "user-123",
		Channel: "sms",
		Message: "Test SMS",
	}

	err := notificationProcessor.Process(notification)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if smsHandler.calls != 1 {
		t.Fatalf(
			"expected handler to be called once, got %d",
			smsHandler.calls,
		)
	}
}

func TestProcessInApp(t *testing.T) {
	inAppHandler := &fakeHandler{}

	notificationProcessor := New(
		map[string]handlers.NotificationHandler{
			"in-app": inAppHandler,
		},
	)

	notification := models.Notification{
		ID:      "test-inapp-001",
		UserID:  "user-123",
		Channel: "in-app",
		Message: "Test in-app notification",
	}

	err := notificationProcessor.Process(notification)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if inAppHandler.calls != 1 {
		t.Fatalf(
			"expected handler to be called once, got %d",
			inAppHandler.calls,
		)
	}
}

func TestProcessUnsupportedChannel(t *testing.T) {
	notificationProcessor := New(
		map[string]handlers.NotificationHandler{},
	)

	notification := models.Notification{
		ID:      "test-invalid-001",
		UserID:  "user-123",
		Channel: "telegram",
		Message: "This should fail",
	}

	err := notificationProcessor.Process(notification)

	if err == nil {
		t.Fatal("expected error for unsupported channel")
	}
}

func TestProcessHandlerError(t *testing.T) {
	expectedErr := errors.New("provider unavailable")

	emailHandler := &fakeHandler{
		err: expectedErr,
	}

	notificationProcessor := New(
		map[string]handlers.NotificationHandler{
			"email": emailHandler,
		},
	)

	notification := models.Notification{
		ID:      "test-provider-error",
		UserID:  "user-123",
		Channel: "email",
		Message: "Provider failure test",
	}

	err := notificationProcessor.Process(notification)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected provider error, got %v",
			err,
		)
	}

	if emailHandler.calls != 1 {
		t.Fatalf(
			"expected handler to be called once, got %d",
			emailHandler.calls,
		)
	}
}
