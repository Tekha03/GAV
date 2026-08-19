package memory

import (
	"context"
	"messenger/internal/model"
	apperrors "shared/app_errors"
	"sync"

	"github.com/google/uuid"
)

type ChatRepository struct {
	mu    sync.RWMutex
	chats map[uuid.UUID]*model.Chat
}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{chats: make(map[uuid.UUID]*model.Chat)}
}

func (cr *ChatRepository) Create(ctx context.Context, chat *model.Chat) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if chat.ID != uuid.Nil {
		if _, found := cr.chats[chat.ID]; found {
			return apperrors.New(apperrors.ChatAlreadyExists, "chat already exists")
		}
	} else {
		chat.ID = uuid.New()
	}

	if _, found := cr.chats[chat.ID]; found {
		return apperrors.New(apperrors.ChatAlreadyExists, "chat already exists")
	}

	cr.chats[chat.ID] = chat
	return nil
}

func (cr *ChatRepository) UpdateTitle(ctx context.Context, chatID uuid.UUID, title string) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	model, found := cr.chats[chatID]
	if !found {
		return apperrors.New(apperrors.ChatNotFound, "chat not found")
	}

	model.Title = title
	return nil
}

func (cr *ChatRepository) UpdatePhoto(ctx context.Context, chatID uuid.UUID, photoURL string) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	model, found := cr.chats[chatID]
	if !found {
		return apperrors.New(apperrors.ChatNotFound, "chat not found")
	}

	model.PhotoURL = photoURL
	return nil
}

func (cr *ChatRepository) Delete(ctx context.Context, chatID uuid.UUID) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if _, found := cr.chats[chatID]; !found {
		return apperrors.New(apperrors.ChatNotFound, "chat not found")
	}

	delete(cr.chats, chatID)
	return nil
}

func (cr *ChatRepository) GetByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	model, found := cr.chats[chatID]
	if !found {
		return nil, apperrors.New(apperrors.ChatNotFound, "chat not found")
	}

	return model, nil
}
