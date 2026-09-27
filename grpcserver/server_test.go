package grpcserver

import (
	"context"
	"errors"
	"testing"

	"notification-agent/proto"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeNotificationStore struct {
	status  string
	channel string
	userID  string
	err     error
}

func (f *fakeNotificationStore) GetNotificationStatus(
	ctx context.Context,
	notificationID string,
) (string, string, string, error) {
	return f.status, f.channel, f.userID, f.err
}

func TestGetNotificationStatusSuccess(t *testing.T) {
	store := &fakeNotificationStore{
		status:  "processed",
		channel: "email",
		userID:  "user-123",
	}

	server := NewServer(store)

	response, err := server.GetNotificationStatus(
		context.Background(),
		&proto.GetNotificationStatusRequest{
			NotificationId: "notification-001",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.GetNotificationId() != "notification-001" {
		t.Errorf(
			"expected notification ID notification-001, got %s",
			response.GetNotificationId(),
		)
	}

	if response.GetStatus() != "processed" {
		t.Errorf(
			"expected status processed, got %s",
			response.GetStatus(),
		)
	}

	if response.GetChannel() != "email" {
		t.Errorf(
			"expected channel email, got %s",
			response.GetChannel(),
		)
	}

	if response.GetUserId() != "user-123" {
		t.Errorf(
			"expected user ID user-123, got %s",
			response.GetUserId(),
		)
	}
}

func TestGetNotificationStatusInvalidArgument(t *testing.T) {
	store := &fakeNotificationStore{}
	server := NewServer(store)

	_, err := server.GetNotificationStatus(
		context.Background(),
		&proto.GetNotificationStatusRequest{},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf(
			"expected InvalidArgument, got %s",
			status.Code(err),
		)
	}
}

func TestGetNotificationStatusNotFound(t *testing.T) {
	store := &fakeNotificationStore{
		err: pgx.ErrNoRows,
	}

	server := NewServer(store)

	_, err := server.GetNotificationStatus(
		context.Background(),
		&proto.GetNotificationStatusRequest{
			NotificationId: "missing-001",
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if status.Code(err) != codes.NotFound {
		t.Fatalf(
			"expected NotFound, got %s",
			status.Code(err),
		)
	}
}

func TestGetNotificationStatusInternalError(t *testing.T) {
	store := &fakeNotificationStore{
		err: errors.New("database connection failed"),
	}

	server := NewServer(store)

	_, err := server.GetNotificationStatus(
		context.Background(),
		&proto.GetNotificationStatusRequest{
			NotificationId: "notification-002",
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if status.Code(err) != codes.Internal {
		t.Fatalf(
			"expected Internal, got %s",
			status.Code(err),
		)
	}
}
