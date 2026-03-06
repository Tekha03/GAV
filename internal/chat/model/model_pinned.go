package model

import (
	"time"

	"github.com/google/uuid"
)

type PinnedMessages struct {
	ChatID 		uuid.UUID
	MessageID 	uuid.UUID
	PinnedAt	*time.Time
}