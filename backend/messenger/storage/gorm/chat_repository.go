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

type ChatRepository struct {
	repo *Repository
}

func NewChatRepository(repo *Repository) repository.ChatRepository {
	return &ChatRepository{repo: repo}
}

func (cr *ChatRepository) Create(ctx context.Context, chat *model.Chat) error {
	err := cr.repo.WithContext(ctx).Create(chat).Error
	return createError(err, apperrors.ChatAlreadyExists, "chat already exists", "failed to create chat")
}

func (cr *ChatRepository) UpdateTitle(ctx context.Context, chatID uuid.UUID, title string) error {
	result := cr.repo.WithContext(ctx).
		Model(&model.Chat{}).
		Where("id = ?", chatID).
		Update("title", title)
	return mutationError(result, apperrors.ChatNotFound, "chat not found", "failed to update chat title")
}

func (cr *ChatRepository) UpdatePhoto(ctx context.Context, chatID uuid.UUID, photoURL string) error {
	result := cr.repo.WithContext(ctx).
		Model(&model.Chat{}).
		Where("id = ?", chatID).
		Update("photo_url", photoURL)
	return mutationError(result, apperrors.ChatNotFound, "chat not found", "failed to update chat photo")
}

func (cr *ChatRepository) Delete(ctx context.Context, chatID uuid.UUID) error {
	result := cr.repo.WithContext(ctx).Delete(&model.Chat{}, "id = ?", chatID)
	return mutationError(result, apperrors.ChatNotFound, "chat not found", "failed to delete chat")
}

func (cr *ChatRepository) GetByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	var chat model.Chat
	err := cr.repo.WithContext(ctx).First(&chat, "id = ?", chatID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ChatNotFound, "chat not found")
		}
		return nil, internalError("failed to get chat by ID", err)
	}
	return &chat, nil
}

func (cr *ChatRepository) GetByPrivateKey(ctx context.Context, privateKey string) (*model.Chat, error) {
	var chat model.Chat
	err := cr.repo.WithContext(ctx).First(&chat, "private_key = ?", privateKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, internalError("failed to get private chat", err)
	}
	return &chat, nil
}
