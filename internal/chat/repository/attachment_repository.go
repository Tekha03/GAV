package repository

import (
	"context"
	"gav/internal/chat"
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *chat.Attachment) error
	GetByMessage(ctx context.Context, messageID uint) (*chat.Attachment, error)
	DeleteByMessage(ctx context.Context, messageID uint) error
}
