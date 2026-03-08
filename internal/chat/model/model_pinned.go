package model

import (
	"time"

	"github.com/google/uuid"
)

type PinnedMessages struct {
	ChatID 		uuid.UUID `gorm:"type:uuid;primaryKey"`
	MessageID 	uuid.UUID `gorm:"type:uuid"`
	PinnedAt	time.Time `gorm:"autoCreateTime"`
}