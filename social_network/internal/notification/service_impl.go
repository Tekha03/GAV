package notification

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type service struct {
	hub *Hub
}

func NewService(hub *Hub) (NotificationService, error) {
	if hub == nil {
		return nil, ErrEmptyHub
	}
	return &service{hub: hub}, nil
}

func (s *service) NotifyLike(ctx context.Context, postOwnerID, likerID, postID uuid.UUID) error {
	notification := NewNotification(postOwnerID, likerID, TypeLike, postID, LikeMessage)
	return s.sendToUser(notification, postOwnerID)
}

func (s *service) NotifyComment(ctx context.Context, postOwnerID, commenterID, postID uuid.UUID) error {
	notification := NewNotification(postOwnerID, commenterID, TypeComment, postID, CommentMessage)
	return s.sendToUser(notification, postOwnerID)
}

func (s *service) NotifyFollow(ctx context.Context, followingID, followerID uuid.UUID) error {
	notification := NewNotification(followingID, followerID, TypeFollow, followerID, FollowMessage)
	return s.sendToUser(notification, followingID)
}

func (s *service) sendToUser(notification *Notification, userID uuid.UUID) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return ErrFailedToMarshal
	}

	s.hub.SendToUser(userID, data)
	return nil
}