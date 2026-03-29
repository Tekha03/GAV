package handlers

import (
	"context"
	"social_network/internal/post"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) Create(ctx context.Context, userID uuid.UUID, content, imageUrl string) (*post.Post, error) {
	args := m.Called(ctx, userID, content, imageUrl)
	return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockPostService) GetByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockPostService) ListByUser(ctx context.Context, userID uuid.UUID) ([]*post.Post, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*post.Post), args.Error(1)
}

func (m *MockPostService) GetFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*post.Post, time.Time, error) {
	args := m.Called(ctx, userID, before, limit)
	return args.Get(0).([]*post.Post), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockPostService) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	args := m.Called(ctx, userID, postID)
	return args.Error(0)
}
