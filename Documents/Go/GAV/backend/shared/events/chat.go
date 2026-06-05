package events

import "github.com/google/uuid"

type ChatCreatedData struct {
	ChatID   uuid.UUID 	`json:"chat_id"`
	Members []uuid.UUID `json:"members"`
}

type ChatMemberAddedData struct {
	ChatID uuid.UUID `json:"chat_id"`
	UserID uuid.UUID `json:"user_id"`
}

type ChatMemberRemovedData struct {
	ChatID uuid.UUID `json:"chat_id"`
	UserID uuid.UUID `json:"user_id"`
}

type ChatDeletedData struct {
	ChatID uuid.UUID `json:"chat_id"`
}
