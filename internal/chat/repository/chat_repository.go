package repository

import (
	"context"
	"gav/internal/chat"

	"github.com/google/uuid"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *chat.Chat) error
	UpdateTitle(ctx context.Context, chatID uuid.UUID, title string) error
	UpdatePhoto(ctx context.Context, chatID uuid.UUID, photoURL string) error
	Delete(ctx context.Context, chatID uuid.UUID) error
	GetByID(ctx context.Context, chatID uuid.UUID) (*chat.Chat, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*chat.Chat, error)
}
