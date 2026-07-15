package notification

import (
	"context"
	"encoding/json"
	"social_network/internal/device"
	"social_network/internal/firebase"

	"github.com/google/uuid"
)

type service struct {
	hub              *Hub
	notificationRepo Repository
	deviceRepo       device.Repository
	firebaseClient   *firebase.Client
}

func NewService(hub *Hub, notificationRepo Repository, deviceRepo device.Repository, firebaseClient *firebase.Client) (NotificationService, error) {
	if hub == nil || firebaseClient == nil {
		return nil, ErrFirebaseClientEmpty
	}

	return &service{
		hub:              hub,
		notificationRepo: notificationRepo,
		deviceRepo:       deviceRepo,
		firebaseClient:   firebaseClient,
	}, nil
}

func (s *service) NotifyLike(ctx context.Context, postOwnerID, likerID, postID uuid.UUID) error {
	return s.notify(ctx, postOwnerID, likerID, TypeLike, postID, LikeMessage, LikeMessage, map[string]string{
		"entity_type": "post",
		"entity_id":   postID.String(),
		"type":        string(TypeLike),
	})
}

func (s *service) NotifyComment(ctx context.Context, postOwnerID, commenterID, postID uuid.UUID) error {
	return s.notify(ctx, postOwnerID, commenterID, TypeComment, postID, CommentMessage, CommentMessage, map[string]string{
		"entity_type": "comment",
		"entity_id":   postID.String(),
		"type":        string(TypeComment),
	})
}

func (s *service) NotifyFollow(ctx context.Context, followingID, followerID uuid.UUID) error {
	return s.notify(ctx, followingID, followerID, TypeFollow, uuid.Nil, FollowMessage, FollowMessage, map[string]string{
		"entity_type": "following",
		"entity_id":   followingID.String(),
		"type":        string(TypeFollow),
	})
}

func (s *service) NotifyNewMessage(ctx context.Context, receiverID, senderID, chatID uuid.UUID) error {
	return s.notify(ctx, receiverID, senderID, TypeDirectMessage, chatID, DirectMessage, DirectMessage, map[string]string{
		"entity_type": "chat",
		"entity_id":   chatID.String(),
		"type":        string(TypeDirectMessage),
	})
}

func (s *service) NotifyChatInvite(ctx context.Context, receiverID, chatID uuid.UUID) error {
	return s.notify(ctx, receiverID, uuid.Nil, TypeChatInvite, chatID, ChatInviteMessage, ChatInviteMessage, map[string]string{
		"entity_type": "chat",
		"entity_id":   chatID.String(),
		"type":        string(TypeChatInvite),
	})
}

func (s *service) NotifyMessageReaction(ctx context.Context, receiverID, reactorID, messageID uuid.UUID) error {
	return s.notify(ctx, receiverID, reactorID, TypeMessageReaction, messageID, MessageReaction, MessageReaction, map[string]string{
		"entity_type": "message",
		"entity_id":   messageID.String(),
		"type":        string(TypeMessageReaction),
	})
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

func (s *service) notify(ctx context.Context, userID, fromUserID uuid.UUID, notificationType Type, entityID uuid.UUID, title, body string, pushData map[string]string) error {
	notification := NewNotification(userID, fromUserID, notificationType, entityID, body)

	if s.notificationRepo != nil {
		if err := s.notificationRepo.Create(ctx, notification); err != nil {
			return err
		}
	}

	if s.hub != nil {
		if err := s.sendToUser(notification, userID); err != nil {
			return err
		}
	}

	if s.deviceRepo != nil && s.firebaseClient != nil {
		tokens, err := s.deviceRepo.GetByUser(ctx, userID)
		if err == nil {
			for _, token := range tokens {
				data := cloneMap(pushData)
				data["notification_id"] = notification.ID.String()
				_ = s.firebaseClient.SendPush(ctx, token.Token, title, body, data)
			}
		}
	}

	return nil
}

func (s *service) sendToUser(notification *Notification, userID uuid.UUID) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return ErrFailedToMarshal
	}

	s.hub.SendToUser(userID, data)
	return nil
}

func cloneMap(input map[string]string) map[string]string {
	result := make(map[string]string, len(input))
	for key, value := range input {
		result[key] = value
	}

	return result
}
