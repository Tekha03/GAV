package handlers

import (
	"context"
	"social_network/internal/post"
	"time"

	"github.com/google/uuid"
)

type MockFeedService struct {
	GetFeedFn func(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*post.Post, time.Time, error)
}

func (m *MockFeedService) GetFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*post.Post, time.Time, error) {
	return m.GetFeedFn(ctx, userID, before, limit)
}
