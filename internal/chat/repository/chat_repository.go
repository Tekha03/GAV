package repository

import (
	"context"
	"gav/internal/chat"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *chat.Chat) error
	Update(ctx context.Context, chat *chat.Chat) error
	Delete(ctx context.Context, chatID uint) error
	GetByID(ctx context.Context, chatID uint) (*chat.Chat, error)
	GetByUserID(ctx context.Context, userID uint) ([]*chat.Chat, error)
}
