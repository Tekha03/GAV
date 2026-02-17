package memory

import (
	"context"
	"errors"
	"gav/internal/chat"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrChatExists = errors.New("chat already exists")
	ErrChatNotFound = errors.New("chat not found")
)

type ChatRepository struct {
	mu 		sync.RWMutex
	chats 	map[uuid.UUID]*chat.Chat
}

func NewChatrepository() *ChatRepository {
	return &ChatRepository{chats: 
		make(map[uuid.UUID]*chat.Chat),
	}
}

func (cr *ChatRepository) Create(ctx context.Context, chat *chat.Chat) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if chat.ID != uuid.Nil {
		if _, found := cr.chats[chat.ID]; found {
			return ErrDogExists
		}
	} else {
		chat.ID = uuid.New()
	}

	if _, found := cr.chats[chat.ID]; found {
		return ErrDogExists
	}

	cr.chats[chat.ID] = chat
	return nil
}

func (cr *ChatRepository) UpdateTitle(ctx context.Context, chatID uuid.UUID, title string) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	chat, found := cr.chats[chatID]
	if !found {
		return ErrChatNotFound
	}

	chat.Title = title
	return nil
}

func (cr *ChatRepository) UpdatePhoto(ctx context.Context, chatID uuid.UUID, photoURL string) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	chat, found := cr.chats[chatID]
	if !found {
		return ErrChatNotFound
	}

	chat.PhotoURL = photoURL
	return nil
}

func (cr *ChatRepository) Delete(ctx context.Context, chatID uuid.UUID) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if _, found := cr.chats[chatID]; !found {
		return ErrChatNotFound
	}

	delete(cr.chats, chatID)
	return nil
}

func (cr *ChatRepository) GetByID(ctx context.Context, chatID uuid.UUID) (*chat.Chat, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	chat, found := cr.chats[chatID]
	if !found {
		return nil, ErrChatNotFound
	}

	return chat, nil
}
