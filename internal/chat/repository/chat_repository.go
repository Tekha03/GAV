package repository

import (
	"context"
	"gav/internal/chat"

	"github.com/google/uuid"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *chat.Chat) error
	Update(ctx context.Context, chat *chat.Chat) error
	Delete(ctx context.Context, chatID uuid.UUID) error
	GetByID(ctx context.Context, chatID uuid.UUID) (*chat.Chat, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*chat.Chat, error)
}
