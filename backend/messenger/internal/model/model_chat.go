package model

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IsGroup   bool      `gorm:"not null;default:false"`
	Title     string    `gorm:"type:text"`
	PhotoURL  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
