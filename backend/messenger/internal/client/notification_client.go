// messanger/chat/client/notification_client.go
package client

import (
	pb "api/gen/notification/v1"
	"context"
	apperrors "shared/app_errors"
	"shared/retry"

	uuid "github.com/google/uuid"
	"google.golang.org/grpc"
)

type NotificationClient struct {
	client pb.NotificationServiceClient
	conn   *grpc.ClientConn
}

func (nc *NotificationClient) Close() error {
	if nc == nil || nc.conn == nil {
		return nil
	}
	return apperrors.Wrap(apperrors.Internal, "failed to close notification client", nc.conn.Close())
}

func NewNotificationClient(addr string) (*NotificationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ServiceUnavailable, "failed to create notification client", err)
	}
	return &NotificationClient{
		client: pb.NewNotificationServiceClient(conn),
		conn:   conn,
	}, nil
}

func (nc *NotificationClient) SendNewMessage(
	ctx context.Context,
	toUserID uuid.UUID,
	senderName, body, chatID string,
) error {
	err := retry.Do(ctx, retry.DefaultConfig(), func(ctx context.Context) error {
		return nc.sendNewMessageOnce(ctx, toUserID, senderName, body, chatID)
	})
	return apperrors.Wrap(apperrors.ServiceUnavailable, "failed to send notification", err)
}

func (nc *NotificationClient) sendNewMessageOnce(
	ctx context.Context,
	toUserID uuid.UUID,
	senderName, body, chatID string,
) error {
	req := &pb.NotificationRequest{
		UserId:    toUserID[:],
		Title:     senderName + " написал вам",
		Body:      body,
		EventType: "new_message",
		Data: map[string]string{
			"chat_id": chatID,
		},
	}
	_, err := nc.client.SendNotification(ctx, req)
	return err
}
