package handlers

import (
	"context"
	"social_network/internal/post"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFeedService struct {
	mock.Mock
}

func (m *MockFeedService) GetFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*post.Post, time.Time, error) {
	args := m.Called(ctx, userID, before, limit)
	return args.Get(0).([]*post.Post), args.Get(1).(time.Time), args.Error(2)
}
