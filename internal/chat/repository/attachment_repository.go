package repository

import (
	"context"
	"gav/internal/chat"

	"github.com/google/uuid"
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *chat.Attachment) error
	GetByMessage(ctx context.Context, messageID uuid.UUID) ([]*chat.Attachment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*chat.Attachment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByMessage(ctx context.Context, messageID uuid.UUID) error
}