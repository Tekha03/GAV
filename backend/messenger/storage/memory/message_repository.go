package memory

import (
	"context"
	"messenger/internal/model"
	"messenger/internal/repository"
	apperrors "shared/app_errors"
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
			return msg.ID, apperrors.New(apperrors.MessageAlreadyExists, "message already exists")
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
		return apperrors.New(apperrors.MessageNotFound, "message not found")
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
		return apperrors.New(apperrors.MessageNotFound, "message not found")
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
		return nil, apperrors.New(apperrors.MessageNotFound, "message not found")
	}

	return msg, nil
}

func (mr *MessageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var result []*model.Message
	var cursorMessage *model.Message
	if cursorID != nil {
		var ok bool
		cursorMessage, ok = mr.messages[*cursorID]
		if !ok || cursorMessage.ChatID != chatID || cursorMessage.DeletedAt != nil {
			return nil, apperrors.New(apperrors.MessageNotFound, "cursor message not found")
		}
	}

	for _, msg := range mr.messages {
		if msg.ChatID != chatID || msg.DeletedAt != nil {
			continue
		}

		if cursorMessage != nil {
			if msg.CreatedAt.After(cursorMessage.CreatedAt) {
				continue
			}
			if msg.CreatedAt.Equal(cursorMessage.CreatedAt) && msg.ID.String() >= cursorMessage.ID.String() {
				continue
			}
		}

		result = append(result, msg)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].ID.String() > result[j].ID.String()
		}
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func (mr *MessageRepository) GetLastMessageIDForChat(ctx context.Context, chatID uuid.UUID) (uuid.UUID, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var last *model.Message
	for _, msg := range mr.messages {
		if msg.ChatID != chatID || msg.DeletedAt != nil {
			continue
		}
		if last == nil ||
			msg.CreatedAt.After(last.CreatedAt) ||
			(msg.CreatedAt.Equal(last.CreatedAt) && msg.ID.String() > last.ID.String()) {
			last = msg
		}
	}

	if last == nil {
		return uuid.Nil, nil
	}

	return last.ID, nil
}

func (mr *MessageRepository) CountUnreadForChat(ctx context.Context, chatID, userID, lastReadMessageID uuid.UUID) (int, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var lastRead *model.Message
	if lastReadMessageID != uuid.Nil {
		var ok bool
		lastRead, ok = mr.messages[lastReadMessageID]
		if !ok || lastRead.ChatID != chatID || lastRead.DeletedAt != nil {
			return 0, apperrors.New(apperrors.MessageNotFound, "last read message not found")
		}
	}

	count := 0
	for _, msg := range mr.messages {
		if msg.ChatID != chatID || msg.DeletedAt != nil || msg.SenderID == userID {
			continue
		}
		if lastRead == nil ||
			msg.CreatedAt.After(lastRead.CreatedAt) ||
			(msg.CreatedAt.Equal(lastRead.CreatedAt) && msg.ID.String() > lastRead.ID.String()) {
			count++
		}
	}

	return count, nil
}
