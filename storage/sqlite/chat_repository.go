package sqlite

import (
	"context"
	"errors"
	"gav/internal/chat/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository struct {
	*BaseRepository
}

func NewChatrepository(db *gorm.DB) (*ChatRepository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &ChatRepository{BaseRepository: repo}, nil
}

func (cr *ChatRepository) Create(ctx context.Context, chat *model.Chat) error {
	if chat.ID == uuid.Nil {
		chat.ID = uuid.New()
	}

	err := cr.DB(ctx).WithContext(ctx).Create(chat).Error
	if err != nil {
		return err
	}

	return nil
}

func (cr *ChatRepository) UpdateTitle(ctx context.Context, chatID uuid.UUID, title string) error {
	if len(title) == 0 {
		return ErrEmptyTitle
	}

	var chat model.Chat
	err := cr.DB(ctx).
		First(&chat, "id = ?", chatID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrChatNotFound
	}
	if err != nil {
		return err
	}

	if !chat.IsGroup {
		return ErrNotGroup
	}

	return cr.DB(ctx).
		Model(&chat).
		Update("title", title).Error
}

func (cr *ChatRepository) UpdatePhoto(ctx context.Context, chatID uuid.UUID, photoURL string) error {
	var chat model.Chat

	err := cr.DB(ctx).
		First(&chat, "id = ?", chatID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrChatNotFound
	}
	if err != nil {
		return err
	}

	if !chat.IsGroup {
		return ErrNotGroup
	}

	return cr.DB(ctx).
		Model(&chat).
		Update("photo_url", photoURL).Error
}

func (cr *ChatRepository) DeletePhoto(ctx context.Context, chatID uuid.UUID) error {
	var chat model.Chat

	err := cr.DB(ctx).
		First(&chat, "id = ?", chatID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrChatNotFound
	}
	if err != nil {
		return err
	}

	if !chat.IsGroup {
		return ErrNotGroup
	}

	return cr.DB(ctx).
		Model(&model.Chat{}).
		Where("id = ?", chatID).
		Update("photo_url", "").Error
}

func (cr *ChatRepository) GetByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	var chat model.Chat

	err := cr.DB(ctx).
		First(&chat, "id = ?", chatID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrChatNotFound
	}
	
	if err != nil {
		return nil, err
	}

	return &chat, nil
}