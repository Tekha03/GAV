package sqlite

import (
	"context"

	"gav/internal/chat/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	*BaseRepository
}

func NewAttachmentRepository(db *gorm.DB) (*AttachmentRepository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &AttachmentRepository{BaseRepository: repo}, nil
}

func (ar *AttachmentRepository) Create(ctx context.Context, attachment *model.Attachment) error {
	if attachment.ID == uuid.Nil {
		attachment.ID = uuid.New()
	}

	return ar.DB(ctx).Create(attachment).Error
}

func (ar *AttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Attachment, error) {
	var att model.Attachment

	err := ar.DB(ctx).
		Where("id = ?", id).
		First(&att).Error

	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (ar *AttachmentRepository) GetByMessage(ctx context.Context, messageID uuid.UUID) ([]*model.Attachment, error) {
	var attachments []*model.Attachment

	err := ar.DB(ctx).
		Where("message_id = ?", messageID).
		Find(&attachments).Error

	if err != nil {
		return nil, err
	}

	return attachments, nil
}

func (ar *AttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return ar.DB(ctx).
		Where("id = ?", id).
		Delete(&model.Attachment{}).Error
}

func (ar *AttachmentRepository) DeleteByMessage(ctx context.Context, messageID uuid.UUID) error {
	return ar.DB(ctx).
		Where("message_id = ?", messageID).
		Delete(&model.Attachment{}).Error
}