package client

import (
	pb "api/gen/notification/v1"
	"context"
	apperrors "shared/app_errors"
	"shared/retry"
	"time"

	uuid "github.com/google/uuid"
	"google.golang.org/grpc"
)

type NotificationClient struct {
	client  pb.NotificationServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

func (nc *NotificationClient) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if nc.timeout <= 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, nc.timeout)
}

func (nc *NotificationClient) Close() error {
	if nc == nil || nc.conn == nil {
		return nil
	}
	return apperrors.Wrap(apperrors.Internal, "failed to close notification client", nc.conn.Close())
}

func NewNotificationClient(addr string) (NotifiClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ServiceUnavailable, "failed to create notification client", err)
	}
	rawClient := &NotificationClient{
		client:  pb.NewNotificationServiceClient(conn),
		conn:    conn,
		timeout: 3 * time.Second,
	}

	return NewProtectedNotificationClient(rawClient), nil
}

func (nc *NotificationClient) SendNewMessage(
	ctx context.Context,
	toUserID, senderID uuid.UUID,
	senderName, body, chatID string,
) error {
	ctx, cancel := nc.withTimeout(ctx)
	defer cancel()

	err := retry.Do(ctx, retry.DefaultConfig(), func(ctx context.Context) error {
		return nc.sendNewMessageOnce(ctx, toUserID, senderID, senderName, body, chatID)
	})
	return apperrors.Wrap(apperrors.ServiceUnavailable, "failed to send notification", err)
}

func (nc *NotificationClient) sendNewMessageOnce(
	ctx context.Context,
	toUserID, senderID uuid.UUID,
	senderName, body, chatID string,
) error {
	ctx, cancel := nc.withTimeout(ctx)
	defer cancel()

	req := &pb.NotificationRequest{
		UserId:    toUserID[:],
		Title:     senderName + " написал вам",
		Body:      body,
		EventType: "new_message",
		Data: map[string]string{
			"chat_id":   chatID,
			"sender_id": senderID.String(),
		},
	}
	_, err := nc.client.SendNotification(ctx, req)
	return err
}
