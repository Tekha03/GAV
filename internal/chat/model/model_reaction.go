package model

import "github.com/google/uuid"

type Reaction struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	MessageID uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Emoji     string
}
