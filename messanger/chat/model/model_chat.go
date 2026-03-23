package model

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID   	`gorm:"type:uuid;primaryKey"`
	IsGroup   bool			`gorm:"not null;default:false"`
	Title     string
	PhotoURL  string
	CreatedAt time.Time		`gorm:"autoCreateTime"`
}