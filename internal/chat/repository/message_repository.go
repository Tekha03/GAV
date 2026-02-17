package repository

import (
	"context"
	"gav/internal/chat"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Create(ctx context.Context, message *chat.Message) error
	UpdateText(ctx context.Context, messageID uuid.UUID, newText string) error
	Delete(ctx context.Context, messageID uuid.UUID) error
	GetByID(ctx context.Context, messageID uuid.UUID) (*chat.Message, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*chat.Message, error)
}
