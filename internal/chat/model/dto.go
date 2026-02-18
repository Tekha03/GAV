package model

import "github.com/google/uuid"

type SendMessageInput struct {
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Text   	  *string
	Type      MessageType
	ReplyToID *uuid.UUID
}

type AttachmentInput struct {
	Type     AttachmentType
	URL      string
	FileName string
	FileSize int64
}
