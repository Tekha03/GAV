package repository

import (
	"context"
	"messanger/internal/model"

	"github.com/google/uuid"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) error
	UpdateTitle(ctx context.Context, chatID uuid.UUID, title string) error
	UpdatePhoto(ctx context.Context, chatID uuid.UUID, photoURL string) error
	Delete(ctx context.Context, chatID uuid.UUID) error
	GetByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error)
}
