package handlers

import (
	"testing"

	"notification-agent/models"
)

func TestEmailHandlerRequiresRecipient(t *testing.T) {
	notification := models.Notification{
		ID:      "test-email-handler",
		UserID:  "user-123",
		Channel: "email",
		Subject: "Test Email",
		Message: "Hello from email test",
	}

	handler := EmailHandler{}

	err := handler.Send(notification)

	if err == nil {
		t.Fatal("expected error when recipient email is missing")
	}

	if err.Error() != "recipient email is required" {
		t.Fatalf(
			"expected recipient email error, got %v",
			err,
		)
	}
}

func TestSMSHandlerRequiresRecipient(t *testing.T) {
	notification := models.Notification{
		ID:      "test-sms-handler",
		UserID:  "user-123",
		Channel: "sms",
		Message: "Hello from SMS test",
	}

	handler := SMSHandler{}

	err := handler.Send(notification)

	if err == nil {
		t.Fatal("expected error when recipient phone is missing")
	}

	if err.Error() != "recipient phone is required" {
		t.Fatalf(
			"expected recipient phone error, got %v",
			err,
		)
	}
}

func TestInAppHandlerSend(t *testing.T) {
	notification := models.Notification{
		ID:      "test-inapp-handler",
		UserID:  "user-123",
		Channel: "in-app",
		Message: "Hello from in-app test",
	}

	handler := InAppHandler{}

	err := handler.Send(notification)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
