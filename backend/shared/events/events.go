package events

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeChatCreated   EventType = "ChatCreated"
	EventTypeMessageSent   EventType = "MessageSent"
	EventTypeUserJoined    EventType = "UserJoined"
	EventTypeUserLeft      EventType = "UserLeft"
	EventTypeChatUpdated   EventType = "ChatUpdated"
	EventTypeMessageEdited  EventType = "MessageEdited"
	EventTypeMessageDeleted EventType = "MessageDeleted"
)

type Event struct {
	EventID		uuid.UUID 	`json:"event_id"`
	EventType 	EventType 	`json:"event_type"`
	Timestamp 	time.Time 	`json:"timestamp"`
	Data 		interface{} `json:"data"`
}

type MessageSentData struct {
	MessageID uuid.UUID `json:"message_id"`
	ChatID    uuid.UUID `json:"chat_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Text      string    `json:"text"`
}
