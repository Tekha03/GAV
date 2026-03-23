package repository

import (
	"time"

	"github.com/google/uuid"
)

type TypingRepository interface {
	SetTyping(chatID, userID uuid.UUID)
	GetTypingUsers(chatID uuid.UUID, timeout time.Duration) []uuid.UUID
	Cleanup(timeout time.Duration)
}