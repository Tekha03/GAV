package dto

import "github.com/google/uuid"

type ReactionRequestDTO struct {
	MessageID uuid.UUID `json:"message_id"`
	Emoji     string    `json:"emoji"`
}

type ReactionResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	Emoji     string    `json:"emoji"`
}