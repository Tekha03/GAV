package notification

import (
	"context"
	"encoding/json"
	"social_network/internal/device"

	"github.com/google/uuid"
)

type service struct {
	hub *Hub
	notificationRepo Repository
	deviceRepo     	 device.Repository
}

func NewService(hub *Hub, notificationRepo Repository, fcmClient *fcm.Client,) (NotificationService, error) {
	if hub == nil {
		return nil, ErrEmptyHub
	}
	return &service{
		hub: hub,
		notificationRepo: notificationRepo,
		fcmClient: fcmClient,
	}, nil
}

func (s *service) NotifyLike(
	ctx context.Context,
	postOwnerID, likerID, postID uuid.UUID,
) error {
	notification := NewNotification(
		postOwnerID,
		likerID,
		TypeLike,
		postID,
		LikeMessage,
	)

	if s.notificationRepo != nil {
		if err := s.notificationRepo.Create(ctx, notification); err != nil {
			return err
		}
	}

	if s.hub != nil {
		data, err := json.Marshal(notification)
		if err != nil {
			return ErrFailedToMarshal
		}
		s.hub.SendToUser(postOwnerID, data)
	}

	if s.deviceRepo != nil && s.fcmClient != nil {
		tokens, err := s.deviceRepo.GetByUser(ctx, postOwnerID)
		if err != nil || len(tokens) == 0 {
			return nil
		}

		// Собираем FCM‑сообщение
		msg := &fcm.Message{
			Data: map[string]string{
				"entity_id":   postID.String(),
				"notification_id": notification.ID.String(),
				"type":        "like",
			},
			Notification: &fcm.{
				Title: "Новый лайк",
				Body:  notification.Message,
			},
		}

		for _, token := range tokens {
			_, err := s.fcmClient.Send(ctx, msg, token.Token)
			if err != nil {
				// логгируем, но не ломаем
			}
		}
	}
}


func (s *service) NotifyComment(ctx context.Context, postOwnerID, commenterID, postID uuid.UUID) error {
	notification := NewNotification(
		postOwnerID,
		commenterID,
		TypeComment,
		postID,
		CommentMessage,
	)

	if s.notificationRepo != nil {
		if err := s.notificationRepo.Create(ctx, notification); err != nil {
			return err
		}
	}

	data, err := json.Marshal(notification)
	if err != nil {
		return ErrFailedToMarshal
	}

	s.hub.SendToUser(postOwnerID, data)

	return nil
}

func (s *service) NotifyFollow(ctx context.Context, followingID, followerID uuid.UUID) error {
	notification := NewNotification(
		followingID,
		followerID,
		TypeFollow,
		followerID,
		FollowMessage,
	)

	if s.notificationRepo != nil {
		if err := s.notificationRepo.Create(ctx, notification); err != nil {
			return err
		}
	}

	data, err := json.Marshal(notification)
	if err != nil {
		return ErrFailedToMarshal
	}

	s.hub.SendToUser(followingID, data)

	return nil
}

func (s *service) GetInAppNotifications(ctx context.Context, userID uuid.UUID) ([]*Notification, error) {
	if s.notificationRepo == nil {
		return nil, ErrNotificationRepoEmpty
	}
	return s.notificationRepo.GetByUserID(ctx, userID)
}

func (s *service) MarkInAppNotificationAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	if s.notificationRepo == nil {
		return ErrNotificationRepoEmpty
	}
	return s.notificationRepo.MarkAsRead(ctx, userID, notificationID)
}

func (s *service) sendToUser(notification *Notification, userID uuid.UUID) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return ErrFailedToMarshal
	}

	s.hub.SendToUser(userID, data)
	return nil
}