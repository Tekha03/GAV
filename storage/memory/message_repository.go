package memory

import (
	"bytes"
	"context"
	"errors"
	"gav/internal/chat"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMessageNotFound = errors.New("messsage not found")
)

type MessageRepository struct {
	mu       sync.RWMutex
	messages map[uuid.UUID]*chat.Message
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{messages: map[uuid.UUID]*chat.Message{}}
}

func (mr *MessageRepository) Create(ctx context.Context, msg *chat.Message) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if msg.ID != uuid.Nil {
		if _, found := mr.messages[msg.ID]; found {
			return ErrAttachmentExist
		}
	} else {
		msg.ID = uuid.New()
	}
	msg.CreatedAt = time.Now()

	mr.messages[msg.ID] = msg
	return nil
}

func (mr *MessageRepository) UpdateText(ctx context.Context, msgID uuid.UUID, newText string) error {
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

func (mr *MessageRepository) Delete(ctx context.Context, msgID uuid.UUID) error {
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

func (mr *MessageRepository) GetByID(ctx context.Context, msgID uuid.UUID) (*chat.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	msg, ok := mr.messages[msgID]
	if !ok {
		return nil, ErrMessageNotFound
	}

	return msg, nil
}

func (mr *MessageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*chat.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var result []*chat.Message
	for _, msg := range mr.messages {
		if msg.ChatID != chatID || msg.DeletedAt != nil {
			continue
		}

		if cursorID != nil && bytes.Compare(msg.ID[:], cursorID[:]) >= 0 {
			continue
		}

		result = append(result, msg)
	}

	sort.Slice(result, func(i, j int) bool {
		return bytes.Compare(result[i].ID[:], result[j].ID[:]) > 0
	})

	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}