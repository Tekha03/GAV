package repository

import (
	"context"
	"gav/internal/chat"
)

type MessageRepository interface {
	Create(ctx context.Context, message *chat.Message) error
	UpdateText(ctx context.Context, newText string) error
	Delete(ctx context.Context, messageID uint) error
	GetByID(ctx context.Context, messageID uint) (*chat.Message, error)
	GetByChatID(ctx context.Context, chatID, limit uint, cursorID *uint) ([]*chat.Message, error)
}
