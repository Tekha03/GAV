package repository

import (
	"context"
	"gav/internal/chat"
)

type MessageRepository interface {
	Create(ctx context.Context, message *chat.Message) error
	UpdateText(ctx context.Context, messageID uint, newText string) error
	Delete(ctx context.Context, messageID uint) error
	GetByID(ctx context.Context, messageID uint) (*chat.Message, error)
	GetByChatID(ctx context.Context, chatID uint, limit int, cursorID *uint) ([]*chat.Message, error)
}
