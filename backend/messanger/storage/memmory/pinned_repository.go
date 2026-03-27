package memory

import (
	"context"
	"messanger/internal/repository"
	"sync"
	"time"

	"github.com/google/uuid"
)

type PinnedRepo struct {
    mu      sync.Mutex
    pinned  map[uuid.UUID]map[uuid.UUID]time.Time
}

func NewPinnedRepo() repository.PinnedRepository {
    return &PinnedRepo{pinned: make(map[uuid.UUID]map[uuid.UUID]time.Time)}
}


func (pr *PinnedRepo) Pin(ctx context.Context, chatID, messageID uuid.UUID) error {
    pr.mu.Lock()
    defer pr.mu.Unlock()

    if _, ok := pr.pinned[chatID]; !ok {
        pr.pinned[chatID] = make(map[uuid.UUID]time.Time)
    }

    pr.pinned[chatID][messageID] = time.Now()
    return nil
}

func (pr *PinnedRepo) Unpin(ctx context.Context, chatID, messageID uuid.UUID) error {
    pr.mu.Lock()
    defer pr.mu.Unlock()

    if _, ok := pr.pinned[chatID]; ok {
        delete(pr.pinned[chatID], messageID)
    }
    return nil
}

func (pr *PinnedRepo) GetByChatID(ctx context.Context, chatID uuid.UUID) []uuid.UUID {
    pr.mu.Lock()
    defer pr.mu.Unlock()

    ids := []uuid.UUID{}
    if messages, ok := pr.pinned[chatID]; ok {
        for id := range messages {
            ids = append(ids, id)
        }
    }
    return ids
}