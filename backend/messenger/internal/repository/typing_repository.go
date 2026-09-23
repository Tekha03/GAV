package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TypingRepository interface {
	SetTyping(ctx context.Context, chatID, userID uuid.UUID) error
	GetTypingUsers(ctx context.Context, chatID uuid.UUID, timeout time.Duration) ([]uuid.UUID, error)
	Cleanup(ctx context.Context, timeout time.Duration) error
}
