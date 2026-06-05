package repository

import (
	"context"
	"messenger/internal/model"

	"github.com/google/uuid"
)

type ReactionRepository interface {
	Add(ctx context.Context, reaction *model.Reaction) error
	Remove(ctx context.Context, messageID, userID uuid.UUID) error
}
