package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeMessageSent    EventType = "message_sent"
	EventTypeMessageEdited  EventType = "message_edited"
	EventTypeMessageDeleted EventType = "message_deleted"

	EventTypeChatCreated       EventType = "chat_created"
	EventTypeChatMemberAdded   EventType = "chat_added_user"
	EventTypeChatMemberRemoved EventType = "chat_removed_user"
	EventTypeChatUpdated       EventType = "chat_updated"
	EventTypeChatDeleted       EventType = "chat_deleted"

	EventTypeReactionAdded   EventType = "reaction_added"
	EventTypeReactionRemoved EventType = "reaction_removed"
)

type Event struct {
	EventID   uuid.UUID       `json:"event_id"`
	EventType EventType       `json:"event_type"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}
