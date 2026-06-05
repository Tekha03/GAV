package handlers

import (
	"context"
	"social_network/internal/like"

	"github.com/stretchr/testify/mock"
)

type MockLikeService struct {
	mock.Mock
}

func (m *MockLikeService) Add(ctx context.Context, like like.Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *MockLikeService) Remove(ctx context.Context, like like.Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}
