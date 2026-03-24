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
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    MessageID uuid.UUID      `gorm:"type:uuid;index;not null"`
    URL       string         `gorm:"type:text;not null"`
    Type      AttachmentType `gorm:"type:varchar(20);not null"`
    FileName  string         `gorm:"type:text"`
    FileSize  int64          `gorm:"not null;default:0"`
}
