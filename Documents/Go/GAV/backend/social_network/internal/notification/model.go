package notification

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeLike			Type = "like"
	TypeComment			Type = "comment"
	TypeFollow			Type = "follow"
	TypeDirectMessage 	Type = "direct_message"
	TypeChatInvite 		Type = "chat_invite"
	TypeMessageReaction Type = "message_reaction"
)

type Notification struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
    UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
    FromUserID  uuid.UUID `json:"from_user_id" gorm:"type:uuid"`
    Type        Type      `json:"type" gorm:"type:text"`
    EntityID    uuid.UUID `json:"entity_id" gorm:"type:uuid"`
    Message     string    `json:"message" gorm:"type:text"`
    IsRead      bool      `json:"is_read" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func NewNotification(userID, fromUserID uuid.UUID, notificationType Type, entityID uuid.UUID, message string) *Notification {
	return &Notification{
		ID:			uuid.New(),
		UserID: 	userID,
		FromUserID: fromUserID,
		Type: 		notificationType,
		EntityID: 	entityID,
		Message: 	message,
		IsRead: 	false,
		CreatedAt: 	time.Now(),
	}
}
