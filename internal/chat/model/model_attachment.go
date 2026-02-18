package model

import "github.com/google/uuid"

type AttachmentType string

var (
	AttachmentImage AttachmentType = "image"
	AttachmentVideo AttachmentType = "video"
	AttachmentVoice AttachmentType = "voice"
	AttachmentFile  AttachmentType = "file"
)

type Attachment struct {
	ID        uuid.UUID
	MessageID uuid.UUID
	URL       string
	Type      AttachmentType
	FileName  string
	FileSize  string
}
