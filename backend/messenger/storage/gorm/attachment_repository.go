package gorm

import (
	"context"
	"errors"
	"messenger/internal/model"
	"messenger/internal/repository"
	apperrors "shared/app_errors"

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
	err := ar.repo.WithContext(ctx).Create(attachment).Error
	return createError(err, apperrors.AttachmentAlreadyExists, "attachment already exists", "failed to create attachment")
}

func (ar *AttachmentRepository) CreateBatch(ctx context.Context, attachments []model.Attachment) error {
	if len(attachments) == 0 {
		return nil
	}

	err := ar.repo.WithContext(ctx).Create(&attachments).Error
	return createError(
		err,
		apperrors.AttachmentAlreadyExists,
		"attachment already exists",
		"failed to create attachments",
	)
}

func (ar *AttachmentRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Attachment, error) {
	var attachment model.Attachment

	err := ar.repo.
		WithContext(ctx).
		First(&attachment, "id = ?", id).
		Error

	if err == nil {
		return &attachment, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.New(
			apperrors.AttachmentNotFound,
			"attachment not found",
		)
	}

	return nil, apperrors.Wrap(
		apperrors.Internal,
		"failed to get attachment by ID",
		err,
	)
}

func (ar *AttachmentRepository) GetByMessage(ctx context.Context, messageID uuid.UUID) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	err := ar.repo.WithContext(ctx).
		Where("message_id = ?", messageID).
		Find(&attachments).Error
	if err != nil {
		return nil, internalError("failed to get attachments by message ID", err)
	}
	return attachments, nil
}

func (ar *AttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := ar.repo.WithContext(ctx).Delete(&model.Attachment{}, "id = ?", id)
	return mutationError(result, apperrors.AttachmentNotFound, "attachment not found", "failed to delete attachment")
}

func (ar *AttachmentRepository) DeleteByMessage(ctx context.Context, messageID uuid.UUID) error {
	result := ar.repo.WithContext(ctx).
		Where("message_id = ?", messageID).
		Delete(&model.Attachment{})
	return internalError("failed to delete attachments by message ID", result.Error)
}
