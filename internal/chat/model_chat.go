package chat

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	IsGroup   bool
	Title     string
	PhotoURL  string
	CreatedAt time.Time
}