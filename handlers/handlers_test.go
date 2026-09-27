package handlers

import (
	"testing"

	"notification-agent/models"
)

func TestEmailHandlerSend(t *testing.T) {
	notification := models.Notification{
		ID:      "test-email-handler",
		UserID:  "user-123",
		Channel: "email",
		Subject: "Test Email",
		Message: "Hello from email test",
	}

	handler := EmailHandler{}

	err := handler.Send(notification)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSMSHandlerSend(t *testing.T) {
	notification := models.Notification{
		ID:      "test-sms-handler",
		UserID:  "user-123",
		Channel: "sms",
		Message: "Hello from SMS test",
	}

	handler := SMSHandler{}

	err := handler.Send(notification)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
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
