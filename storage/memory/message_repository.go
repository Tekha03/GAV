package memory

import (
	"context"
	"errors"
	"gav/internal/chat"
	"sort"
	"sync"
	"time"
)

var (
	ErrMessageNotFound = errors.New("messsage not found")
)

type MessageRepository struct {
	mu       sync.RWMutex
	lastID   uint
	messages map[uint]*chat.Message
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{messages: map[uint]*chat.Message{}}
}

func (mr *MessageRepository) Create(ctx context.Context, msg *chat.Message) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	mr.lastID++
	msg.ID = mr.lastID
	msg.CreatedAt = time.Now()

	mr.messages[msg.ID] = msg
	return nil
}

func (mr *MessageRepository) UpdateText(ctx context.Context, msgID uint, newText string) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	msg, ok := mr.messages[msgID]
	if !ok {
		return ErrMessageNotFound
	}

	msg.Content = newText
	now := time.Now()
	msg.EditedAt = &now
	return nil
}

func (mr *MessageRepository) Delete(ctx context.Context, msgID uint) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	msg, ok := mr.messages[msgID]
	if !ok {
		return ErrMessageNotFound
	}

	now := time.Now()
	msg.DeletedAt = &now
	return nil
}

func (mr *MessageRepository) GetByID(ctx context.Context, msgID uint) (*chat.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	msg, ok := mr.messages[msgID]
	if !ok {
		return nil, ErrMessageNotFound
	}

	return msg, nil
}

func (mr *MessageRepository) GetByChatID(ctx context.Context, chatID uint, limit int, cursorID *uint) ([]*chat.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var result []*chat.Message
	for _, msg := range mr.messages {
		if msg.ChatID != chatID || msg.DeletedAt != nil {
			continue
		}

		if cursorID != nil && msg.ID >= *cursorID {
			continue
		}

		result = append(result, msg)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID > result[j].ID
	})

	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}