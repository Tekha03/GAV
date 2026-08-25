package model

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	EventType     string    `gorm:"type:varchar(100);not null"`
	Payload       []byte    `gorm:"type:jsonb;not null"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	PublishedAt   *time.Time
	Attempts      int       `gorm:"not null;default:0"`
	NextAttemptAt time.Time `gorm:"not null"`
	LockedUntil   *time.Time
	LastError     *string
}
