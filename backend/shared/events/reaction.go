package events

import "github.com/google/uuid"

type ReactionAddedData struct {
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	Reaction  string    `json:"reaction"`
}

type ReactionRemovedData struct {
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	Reaction  string    `json:"reaction"`
}
