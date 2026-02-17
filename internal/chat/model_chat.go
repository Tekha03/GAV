package chat

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	IsGroup   bool
	Title     string
	CreatedAt time.Time
}