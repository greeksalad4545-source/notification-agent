package grpcserver

import (
	"context"
	"fmt"

	"notification-agent/proto"
	"notification-agent/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedNotificationServiceServer
	db *storage.Postgres
}

func NewServer(db *storage.Postgres) *Server {
	return &Server{
		db: db,
	}
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

	statusValue, channel, userID, err := s.db.GetNotificationStatus(
		ctx,
		req.GetNotificationId(),
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("failed to get notification status: %v", err),
		)
	}

	return &proto.GetNotificationStatusResponse{
		NotificationId: req.GetNotificationId(),
		Status:         statusValue,
		Channel:        channel,
		UserId:         userID,
	}, nil
}
