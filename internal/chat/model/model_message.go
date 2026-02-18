package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

type Message struct {
	ID       	uuid.UUID
	ChatID   	uuid.UUID
	SenderID 	uuid.UUID

	Text 	 	*string
	ReplyToID 	*uuid.UUID

	CreatedAt	time.Time

	EditedAt  	*time.Time
	DeletedAt 	*time.Time
	ReadAt   	*time.Time
}
