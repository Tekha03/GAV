// messanger/chat/client/notification_client.go
package client

import (
	"context"
	"google.golang.org/grpc"
	pb "api/notification_gen"
	uuid "github.com/google/uuid"
)

type NotificationClient struct {
	client pb.NotificationServiceClient
	conn   *grpc.ClientConn
}

func NewNotificationClient(addr string) (*NotificationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
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
