package user

import (
	"time"

	"github.com/google/uuid"
)

type WalkSession struct {
	ID         uuid.UUID          `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID          `gorm:"type:uuid;not null;index"`
	Lat        float64            `gorm:"not null"`
	Lon        float64            `gorm:"not null"`
	Visibility LocationVisibility `gorm:"not null"`
	StartedAt  time.Time          `gorm:"not null"`
	UpdatedAt  time.Time          `gorm:"not null;index"`
	EndedAt    *time.Time         `gorm:"index"`
}
