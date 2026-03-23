package feed

import (
	"context"
	"social_network/internal/post"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, p *post.Post) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPostRepository) Update(ctx context.Context, p *post.Post) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepository) GetByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockPostRepository) ListFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*post.Post, error) {
	args := m.Called(ctx, userID, before, limit)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*post.Post), args.Error(1)
}

func (m *MockPostRepository) ListByUser(ctx context.Context, authorID uuid.UUID) ([]*post.Post, error) {
	args := m.Called(ctx, authorID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*post.Post), args.Error(1)
}