package grpcserver

import (
	"context"
	"errors"
	"fmt"

	"notification-agent/proto"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationStore interface {
	GetNotificationStatus(
		ctx context.Context,
		notificationID string,
	) (string, string, string, error)
}

type Server struct {
	proto.UnimplementedNotificationServiceServer
	db NotificationStore
}

func NewServer(db NotificationStore) *Server {
	return &Server{db: db}
}

func (s *Server) GetNotificationStatus(
	ctx context.Context,
	req *proto.GetNotificationStatusRequest,
) (*proto.GetNotificationStatusResponse, error) {
	if req.GetNotificationId() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"notification_id is required",
		)
	}

	statusValue, channel, userID, err :=
		s.db.GetNotificationStatus(
			ctx,
			req.GetNotificationId(),
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(
				codes.NotFound,
				"notification not found",
			)
		}

		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf(
				"failed to get notification status: %v",
				err,
			),
		)
	}

	return &proto.GetNotificationStatusResponse{
		NotificationId: req.GetNotificationId(),
		Status:         statusValue,
		Channel:        channel,
		UserId:         userID,
	}, nil
}
