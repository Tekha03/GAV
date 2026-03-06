package memory

import (
	"bytes"
	"context"
	"gav/internal/chat/model"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MessageRepository struct {
	mu       		sync.RWMutex
	messages 		map[uuid.UUID]map[uuid.UUID]*model.Message
	messageReads 	map[uuid.UUID]map[uuid.UUID]time.Time
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{messages: make(map[uuid.UUID]map[uuid.UUID]*model.Message)}
}

func (mr *MessageRepository) Create(ctx context.Context, msg *model.Message) (uuid.UUID, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}
	msg.CreatedAt = time.Now()

	if _, ok := mr.messages[msg.ChatID]; !ok {
		mr.messages[msg.ChatID] = make(map[uuid.UUID]*model.Message)
	}

	if _, exists := mr.messages[msg.ChatID][msg.ID]; exists {
		return uuid.Nil, ErrMessageExists
	}

	mr.messages[msg.ChatID][msg.ID] = msg
	return msg.ID, nil
}

func (mr *MessageRepository) UpdateText(ctx context.Context, chatID, msgID uuid.UUID, newText string) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	chatMessages, ok := mr.messages[chatID]
	if !ok {
		return ErrMessageNotFound
	}

	msg, ok := chatMessages[msgID]
	if !ok {
		return ErrMessageNotFound
	}

	msg.Text = &newText
	now := time.Now()
	msg.EditedAt = &now
	return nil
}

func (mr *MessageRepository) Delete(ctx context.Context, chatID, msgID uuid.UUID) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	chatMessages, ok := mr.messages[chatID]
	if !ok {
		return ErrMessageNotFound
	}

	msg, ok := chatMessages[msgID]
	if !ok {
		return ErrMessageNotFound
	}

	now := time.Now()
	msg.DeletedAt = &now
	return nil
}

func (mr *MessageRepository) GetByID(ctx context.Context, chatID, msgID uuid.UUID) (*model.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	chatMessages, ok := mr.messages[chatID]
	if !ok {
		return nil, ErrMessageNotFound
	}

	msg, ok := chatMessages[msgID]
	if !ok {
		return nil, ErrMessageNotFound
	}

	return msg, nil
}

func (mr *MessageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	chatMessages, ok := mr.messages[chatID]
	if !ok {
		return nil, nil
	}

	var msgs []*model.Message
	for _, msg := range chatMessages {
		if msg.DeletedAt != nil {
			continue
		}

		if cursorID != nil && bytes.Compare(msg.ID[:], cursorID[:]) >= 0 {
			continue
		}

		msgs = append(msgs, msg)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].CreatedAt.Before(msgs[j].CreatedAt)
	})

	if len(msgs) > limit {
		msgs = msgs[:limit]
	}

	return msgs, nil
}

func (mr *MessageRepository) UpdateReadAtForChat(ctx context.Context, chatID, userID uuid.UUID, readAt time.Time) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	chatMessages, ok := mr.messages[chatID]
	if !ok {
		return nil
	}

	for _, msg := range chatMessages {
		if msg.SenderID == userID {
			continue
		}

		if _, ok := mr.messageReads[msg.ID]; !ok {
			mr.messageReads[msg.ID] = make(map[uuid.UUID]time.Time)
		}

		mr.messageReads[msg.ID][userID] = readAt
	}

	return nil
}