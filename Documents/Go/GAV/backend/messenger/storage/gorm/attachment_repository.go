package gorm

import (
	"context"
	"errors"
	"messenger/internal/model"
	"messenger/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
    repo *Repository
}

func NewAttachmentRepository(repo *Repository) repository.AttachmentRepository {
    return &AttachmentRepository{repo: repo}
}

func (ar *AttachmentRepository) Create(ctx context.Context, attachment *model.Attachment) error {
    return ar.repo.WithContext(ctx).Create(attachment).Error
}

func (ar *AttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Attachment, error) {
    var att model.Attachment
    err := ar.repo.WithContext(ctx).First(&att, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrAttachmentNotFound
        }
        return nil, err
    }
    return &att, nil
}

func (ar *AttachmentRepository) GetByMessage(ctx context.Context, messageID uuid.UUID) ([]*model.Attachment, error) {
    var attachments []*model.Attachment
    return attachments, ar.repo.WithContext(ctx).
        Where("message_id = ?", messageID).
        Find(&attachments).Error
}

func (ar *AttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return ar.repo.WithContext(ctx).
        Delete(&model.Attachment{}, "id = ?", id).Error
}

func (ar *AttachmentRepository) DeleteByMessage(ctx context.Context, messageID uuid.UUID) error {
    return ar.repo.WithContext(ctx).
        Where("message_id = ?", messageID).
        Delete(&model.Attachment{}).Error
}
