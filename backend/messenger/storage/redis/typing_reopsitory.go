package redis

import (
	"context"
	"messenger/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TypingRepository struct {
	client *redis.Client
}

func NewTypingRepository(client *redis.Client) repository.TypingRepository {
	return &TypingRepository{client: client}
}

func (tr *TypingRepository) SetTyping(ctx context.Context, chatID, userID uuid.UUID) error {
	key := "typing:" + chatID.String()
	pipe := tr.client.TxPipeline()
	pipe.SAdd(ctx, key, userID.String())
	pipe.Expire(ctx, key, 10*time.Second)
	_, err := pipe.Exec(ctx)
	return unavailable("failed to set typing state", err)
}

func (tr *TypingRepository) GetTypingUsers(ctx context.Context, chatID uuid.UUID, timeout time.Duration) ([]uuid.UUID, error) {
	key := "typing:" + chatID.String()
	usersStr, err := tr.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, unavailable("failed to get typing users", err)
	}

	users := make([]uuid.UUID, 0, len(usersStr))
	for _, userStr := range usersStr {
		userID, err := uuid.Parse(userStr)
		if err != nil {
			return nil, internal("invalid typing user ID in redis", err)
		}
		users = append(users, userID)
	}
	return users, nil
}

func (tr *TypingRepository) Cleanup(ctx context.Context, timeout time.Duration) error {
	return nil
}
