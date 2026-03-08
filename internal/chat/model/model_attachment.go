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
	ID        uuid.UUID		 `gorm:"type:uuid;primaryKey"`
	MessageID uuid.UUID		 `gorm:"type:uuid;index;not null"`
	URL       string		 `gorm:"not null"`
	Type      AttachmentType `gorm:"type:text;not null"`
	FileName  string
	FileSize  int64
}
