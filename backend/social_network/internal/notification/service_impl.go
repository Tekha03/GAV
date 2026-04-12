package notification

import (
	"context"
	"encoding/json"
	"social_network/internal/device"
	"social_network/internal/firebase"

	"github.com/google/uuid"
)

type service struct {
	hub *Hub
	notificationRepo Repository
	deviceRepo     	 device.Repository
	firebaseClient   *firebase.Client
}

func NewService(hub *Hub, notificationRepo Repository, deviceRepo device.Repository, firebaseClient *firebase.Client) (NotificationService, error) {
	if hub == nil || firebaseClient == nil {
		return nil, ErrFirebaseClientEmpty
	}

	return &service{
		hub: hub,
		notificationRepo: notificationRepo,
		deviceRepo: deviceRepo,
		firebaseClient: firebaseClient,
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

	if s.deviceRepo != nil && s.firebaseClient != nil {
		tokens, err := s.deviceRepo.GetByUser(ctx, postOwnerID)
		if err != nil || len(tokens) == 0 {
			return nil
		}

		for _, token := range tokens {

			data := map[string]string{
				"entity_type": "post",
				"entity_id":   postID.String(),
				"notification_id": notification.ID.String(),
				"type":        string(TypeLike),
			}

			err := s.firebaseClient.SendPush(
				ctx,
				token.Token,
				LikeMessage,
				notification.Message,
				data,
			)
			if err != nil {
			}
		}
	}

	return nil
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

	if s.hub != nil {
		data, err := json.Marshal(notification)
		if err != nil {
			return ErrFailedToMarshal
		}
		s.hub.SendToUser(postOwnerID, data)
	}

	if s.deviceRepo != nil && s.firebaseClient != nil {
		tokens, err := s.deviceRepo.GetByUser(ctx, postOwnerID)
		if err != nil || len(tokens) == 0 {
			return nil
		}

		for _, token := range tokens {
			data := map[string]string{
				"entity_type": "comment",
				"entity_id":   postID.String(),
				"notification_id": notification.ID.String(),
				"type":        string(TypeLike),
			}

			err := s.firebaseClient.SendPush(
				ctx,
				token.Token,
				CommentMessage,
				notification.Message,
				data,
			)
			if err != nil {
			}
		}
	}

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

	if s.hub != nil {
		data, err := json.Marshal(notification)
		if err != nil {
			return ErrFailedToMarshal
		}
		s.hub.SendToUser(followingID, data)
	}

	if s.deviceRepo != nil && s.firebaseClient != nil {
		tokens, err := s.deviceRepo.GetByUser(ctx, followingID)
		if err != nil || len(tokens) == 0 {
			return nil
		}

		for _, token := range tokens {
			data := map[string]string{
				"entity_type": "following",
				"entity_id":   followingID.String(),
				"notification_id": notification.ID.String(),
				"type":        string(TypeFollow),
			}

			err := s.firebaseClient.SendPush(
				ctx,
				token.Token,
				FollowMessage,
				notification.Message,
				data,
			)
			if err != nil {
			}
		}
	}

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