package model

import "github.com/google/uuid"

type SendMessageInput struct {
    ChatID      uuid.UUID
    SenderID    uuid.UUID
    Text        *string
    ReplyToID   *uuid.UUID
    Attachments []AttachmentInput
}

type AttachmentInput struct {
	Type     AttachmentType
	URL      string
	FileName string
	FileSize int64
}
