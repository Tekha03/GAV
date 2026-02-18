package repository

import (
	"context"
	"gav/internal/chat/model"

	"github.com/google/uuid"
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *model.Attachment) error
	GetByMessage(ctx context.Context, messageID uuid.UUID) ([]*model.Attachment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Attachment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByMessage(ctx context.Context, messageID uuid.UUID) error
}