package chat

import "time"

type MessageType string

const (
	MessageTypeImage MessageType = "image"
	MessageTypeText  MessageType = "text"
	MessageTypeFile MessageType = "file"
	MessageTypeVideo MessageType = "video"
	MessageTypeSystem MessageType = "system"
	MessageTypeVoice MessageType = "voice"
)

type Message struct {
	ID			uint
	ChatID		uint
	SenderID	uint

	Content		string
	Type		MessageType

	ReplyToID	*uint

	CreatedAt	time.Time

	EditedAt	*time.Time
	DeletedAt	*time.Time
	ReadAt		*time.Time
}