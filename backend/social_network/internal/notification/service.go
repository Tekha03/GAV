package notification

import (
	"context"

	"github.com/google/uuid"
)

type NotificationService interface {
	NotifyLike(ctx context.Context, postOwnerID, likerID, postID uuid.UUID) error
	NotifyComment(ctx context.Context, postOwnerID, commenterID, postID uuid.UUID) error
	NotifyFollow(ctx context.Context, followingID, followerID uuid.UUID) error
	GetInAppNotifications(ctx context.Context, userID uuid.UUID) ([]*Notification, error)
	MarkInAppNotificationAsRead(ctx context.Context, userID, notificationID uuid.UUID) error
}