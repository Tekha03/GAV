package repository

import (
	"context"

	"github.com/google/uuid"
)

type PinnedRepository interface {
	Pin(ctx context.Context, chatID, messageID uuid.UUID) error
	Unpin(ctx context.Context, chatID, messageID uuid.UUID) error
	GetByChatID(ctx context.Context, chatID uuid.UUID) []uuid.UUID
}