package model

import (
	"time"

	"github.com/google/uuid"
)

type MemberRole string

const (
	Admin  MemberRole = "admin"
	Member MemberRole = "member"
)

type ChatMember struct {
    ChatID            uuid.UUID    `gorm:"type:uuid;primaryKey"`
    UserID            uuid.UUID    `gorm:"type:uuid;primaryKey"`
    Role              MemberRole   `gorm:"type:varchar(20);not null;default:'member'"`
    JoinedAt          time.Time    `gorm:"autoCreateTime"`
    LeftAt            *time.Time
    Muted             bool         `gorm:"not null;default:false"`
    LastReadMessageID uuid.UUID    `gorm:"type:uuid"`
}
