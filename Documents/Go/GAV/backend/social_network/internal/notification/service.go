package notification

import (
	"context"

	"github.com/google/uuid"
)

type NotificationService interface {
	NotifyLike(ctx context.Context, postOwnerID, likerID, postID uuid.UUID) error
	NotifyComment(ctx context.Context, postOwnerID, commenterID, postID uuid.UUID) error
	NotifyFollow(ctx context.Context, followingID, followerID uuid.UUID) error
	NotifyNewMessage(ctx context.Context, receiverID, senderID, chatID uuid.UUID) error
	NotifyChatInvite(ctx context.Context, receiverID, chatID uuid.UUID) error
	NotifyMessageReaction(ctx context.Context, receiverID, reactorID, messageID uuid.UUID) error
	GetInAppNotifications(ctx context.Context, userID uuid.UUID) ([]*Notification, error)
	MarkInAppNotificationAsRead(ctx context.Context, userID, notificationID uuid.UUID) error
}
