package repository

import (
	"context"
	"gav/internal/chat"
)

type ReactionRepository interface {
	Add(ctx context.Context, reaction *chat.Reaction) error
	Remove(ctx context.Context, messageID, userID uint) error
	ListByMessage(ctx context.Context, messageID uint) ([]*chat.Reaction, error)
}
