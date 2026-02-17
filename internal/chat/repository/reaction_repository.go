package repository

import (
	"context"
	"gav/internal/chat"

	"github.com/google/uuid"
)

type ReactionRepository interface {
	Add(ctx context.Context, reaction *chat.Reaction) error
	Remove(ctx context.Context, messageID, userID uuid.UUID) error
}
