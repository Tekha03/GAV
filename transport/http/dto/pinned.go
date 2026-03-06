package dto

import (
	"time"

	"github.com/google/uuid"
)

type PinnedMessagesRequestDTO struct {
	ChatID    uuid.UUID `json:"chat_id"`
	MessageID uuid.UUID `json:"message_id"`
}

type PinnedMessagesResponseDTO struct {
	ChatID    uuid.UUID `json:"chat_id"`
	MessageID uuid.UUID `json:"message_id"`
	PinnedAt  time.Time `json:"pinned_at"`
}