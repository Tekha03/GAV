package gorm

import (
	"context"
	"errors"
	"messanger/chat/model"
	"messanger/chat/repository"

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
    return cr.repo.WithContext(ctx).Create(chat).Error
}

func (cr *ChatRepository) UpdateTitle(ctx context.Context, chatID uuid.UUID, title string) error {
    return cr.repo.WithContext(ctx).
        Model(&model.Chat{}).
        Where("id = ?", chatID).
        Update("title", title).Error
}

func (cr *ChatRepository) UpdatePhoto(ctx context.Context, chatID uuid.UUID, photoURL string) error {
    return cr.repo.WithContext(ctx).
        Model(&model.Chat{}).
        Where("id = ?", chatID).
        Update("photo_url", photoURL).Error
}

func (cr *ChatRepository) Delete(ctx context.Context, chatID uuid.UUID) error {
    return cr.repo.WithContext(ctx).
        Delete(&model.Chat{}, "id = ?", chatID).Error
}

func (cr *ChatRepository) GetByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
    var chat model.Chat
    err := cr.repo.WithContext(ctx).First(&chat, "id = ?", chatID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrChatNotFound
        }
        return nil, err
    }
    return &chat, nil
}