package repository

import (
	"context"
	"messenger/internal/model"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) (uuid.UUID, error)
	UpdateText(ctx context.Context, messageID uuid.UUID, newText string) error
	Delete(ctx context.Context, messageID uuid.UUID) error
	GetByID(ctx context.Context, messageID uuid.UUID) (*model.Message, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error)
	UpdateLastReadMessageForChat(ctx context.Context, chatID, userID uuid.UUID) error
}
