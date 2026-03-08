package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

type Message struct {
	ID       	uuid.UUID 	`gorm:"type:uuid;primaryKey"`
	ChatID   	uuid.UUID 	`gorm:"type:uuid;index;not null"`
	SenderID 	uuid.UUID	`gorm:"type:uuid;not null"`

	Text 	 	*string
	ReplyToID 	*uuid.UUID

	CreatedAt	time.Time	`gorm:"autoCreateTime"`

	EditedAt  	*time.Time
	DeletedAt 	*time.Time
	ReadAt   	*time.Time
}
