package redis

import (
	"context"
	"messenger/internal/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PinnedRepository struct {
	client *redis.Client
}

func NewPinnedRepository(client *redis.Client) repository.PinnedRepository {
	return &PinnedRepository{client: client}
}

func (pr *PinnedRepository) Pin(ctx context.Context, chatID, messageID uuid.UUID) error {
	key := "pinned:" + chatID.String()
	return unavailable("failed to pin message", pr.client.LPush(ctx, key, messageID.String()).Err())
}

func (pr *PinnedRepository) Unpin(ctx context.Context, chatID, messageID uuid.UUID) error {
	key := "pinned:" + chatID.String()
	return unavailable("failed to unpin message", pr.client.LRem(ctx, key, 1, messageID.String()).Err())
}

func (pr *PinnedRepository) GetByChatID(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	key := "pinned:" + chatID.String()
	idsStr, err := pr.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, unavailable("failed to get pinned messages", err)
	}

	result := make([]uuid.UUID, 0, len(idsStr))
	for _, idStr := range idsStr {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, internal("invalid pinned message ID in redis", err)
		}
		result = append(result, id)
	}
	return result, nil
}
