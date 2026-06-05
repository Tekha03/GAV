package events

import "github.com/google/uuid"

type MessageSentData struct {
	MessageID uuid.UUID `json:"message_id"`
	ChatID    uuid.UUID `json:"chat_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Text      string    `json:"text"`
}

type MessageEditedData struct {
	MessageID uuid.UUID `json:"message_id"`
	ChatID    uuid.UUID `json:"chat_id"`
	Text      string    `json:"new_text"`
}

type MessageDeletedData struct {
	MessageID uuid.UUID `json:"message_id"`
	ChatID    uuid.UUID `json:"chat_id"`
}
