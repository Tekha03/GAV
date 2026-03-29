package redis

import (
    "context"
    "github.com/redis/go-redis/v9"
    "github.com/google/uuid"
    "messanger/internal/repository"
)

type PinnedRepository struct {
    client *redis.Client
}

func NewPinnedRepository(client *redis.Client) repository.PinnedRepository {
    return &PinnedRepository{client: client}
}

func (pr *PinnedRepository) Pin(ctx context.Context, chatID, messageID uuid.UUID) error {
    key := "pinned:" + chatID.String()
    return pr.client.LPush(ctx, key, messageID.String()).Err()
}

func (pr *PinnedRepository) Unpin(ctx context.Context, chatID, messageID uuid.UUID) error {
    key := "pinned:" + chatID.String()
    return pr.client.LRem(ctx, key, 1, messageID.String()).Err()
}

func (pr *PinnedRepository) GetByChatID(ctx context.Context, chatID uuid.UUID) []uuid.UUID {
    key := "pinned:" + chatID.String()
    idsStr, err := pr.client.LRange(ctx, key, 0, -1).Result()
    if err != nil {
        return nil
    }
    
    result := make([]uuid.UUID, len(idsStr))
    for i, idStr := range idsStr {
        result[i] = uuid.MustParse(idStr)
    }
    return result
}
