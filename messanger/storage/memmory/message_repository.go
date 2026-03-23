package memory

import (
	"bytes"
	"context"
	"messanger/chat/model"
	"messanger/chat/repository"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)


type MessageRepository struct {
	mu       sync.RWMutex
	messages map[uuid.UUID]*model.Message
}

func NewMessageRepository() repository.MessageRepository {
	return &MessageRepository{messages: map[uuid.UUID]*model.Message{}}
}

func (mr *MessageRepository) Create(ctx context.Context, msg *model.Message) (uuid.UUID, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if msg.ID != uuid.Nil {
		if _, found := mr.messages[msg.ID]; found {
			return msg.ID, repository.ErrMessageExists
		}
	} else {
		msg.ID = uuid.New()
	}
	msg.CreatedAt = time.Now()

	mr.messages[msg.ID] = msg
	return msg.ID, nil
}

func (mr *MessageRepository) UpdateText(ctx context.Context, msgID uuid.UUID, newText string) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	msg, ok := mr.messages[msgID]
	if !ok {
		return repository.ErrMessageNotFound
	}

	msg.Text = &newText
	now := time.Now()
	msg.EditedAt = &now
	return nil
}

func (mr *MessageRepository) Delete(ctx context.Context, msgID uuid.UUID) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	msg, ok := mr.messages[msgID]
	if !ok {
		return repository.ErrMessageNotFound
	}

	now := time.Now()
	msg.DeletedAt = &now
	return nil
}

func (mr *MessageRepository) GetByID(ctx context.Context, msgID uuid.UUID) (*model.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	msg, ok := mr.messages[msgID]
	if !ok {
		return nil, repository.ErrMessageNotFound
	}

	return msg, nil
}

func (mr *MessageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var result []*model.Message
	for _, msg := range mr.messages {
		if msg.ChatID != chatID || msg.DeletedAt != nil {
			continue
		}

		if cursorID != nil && msg.ID.String() <= cursorID.String() {
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
func (mr *MessageRepository) UpdateReadAtForChat(ctx context.Context, chatID, userID uuid.UUID, readAt time.Time) error {
    mr.mu.Lock()
    defer mr.mu.Unlock()
    return nil 
}