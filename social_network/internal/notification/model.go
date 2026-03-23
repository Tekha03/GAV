package notification

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeLike	Type = "like"
	TypeComment	Type = "comment"
	TypeFollow	Type = "follow"
)

type Notification struct {
	ID			uuid.UUID	`json:"id"`
	UserID		uuid.UUID	`json:"user_id"`		// кому пришло
	FromUserID	uuid.UUID	`json:"from_user_id"`	// от кого
	Type		Type		`json:"type"`
	EntityID	uuid.UUID	`json:"entity_id"`		// ID поста/комментария
	Message		string		`json:"message"`
	IsRead		bool		`json:"is_read"`
	CreatedAt	time.Time	`json:"created_at"`
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