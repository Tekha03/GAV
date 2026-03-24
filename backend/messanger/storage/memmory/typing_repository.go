package memory

import (
	"context"
	"messanger/chat/repository"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TypingRepository struct {
    mu     sync.Mutex
    typing map[uuid.UUID]map[uuid.UUID]time.Time
}

func NewTypingRepository() repository.TypingRepository {
    return &TypingRepository{
        typing: make(map[uuid.UUID]map[uuid.UUID]time.Time),
    }
}

func (tr *TypingRepository) SetTyping(ctx context.Context, chatID, userID uuid.UUID) {
    tr.mu.Lock()
    defer tr.mu.Unlock()

    if _, ok := tr.typing[chatID]; !ok {
        tr.typing[chatID] = make(map[uuid.UUID]time.Time)
    }

    tr.typing[chatID][userID] = time.Now()
}

func (tr *TypingRepository) GetTypingUsers(ctx context.Context, chatID uuid.UUID, timeout time.Duration) []uuid.UUID {
    tr.mu.Lock()
    defer tr.mu.Unlock()

    now := time.Now()
    users := []uuid.UUID{}

    if chatTyping, ok := tr.typing[chatID]; ok {
        for userID, last := range chatTyping {
            if now.Sub(last) <= timeout {
                users = append(users, userID)
            }
        }
    }

    return users
}

func (tr *TypingRepository) Cleanup(ctx context.Context, timeout time.Duration) {
    tr.mu.Lock()
    defer tr.mu.Unlock()

    now := time.Now()
    for _, users := range tr.typing {
        for userID, last := range users {
            if now.Sub(last) > timeout {
                delete(users, userID)
            }
        }
    }
}