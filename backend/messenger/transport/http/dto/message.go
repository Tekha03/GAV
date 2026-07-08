package dto

import (
	"github.com/google/uuid"
	"time"
)

type MessageResponseDTO struct {
	ID       uuid.UUID `json:"id"`
	ChatID   uuid.UUID `json:"chat_id"`
	SenderID uuid.UUID `json:"sender_id"`

	Text      *string    `json:"text,omitempty"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	EditedAt  *time.Time `json:"edited_at,omitempty"`
}
