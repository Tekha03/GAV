package redis

import (
	"context"
	"messanger/chat/repository"
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

func (tr *TypingRepository) SetTyping(ctx context.Context, chatID, userID uuid.UUID) {
    key := "typing:" + chatID.String()
    tr.client.SAdd(ctx, key, userID.String())
    tr.client.Expire(ctx, key, 10*time.Second)
}

func (tr *TypingRepository) GetTypingUsers(ctx context.Context, chatID uuid.UUID, timeout time.Duration) []uuid.UUID {
    key := "typing:" + chatID.String()
    usersStr, _ := tr.client.SMembers(ctx, key).Result()
    
    users := make([]uuid.UUID, len(usersStr))
    for i, userStr := range usersStr {
        users[i] = uuid.MustParse(userStr)
    }
    return users
}

func (tr *TypingRepository) Cleanup(ctx context.Context, timeout time.Duration) {
}
