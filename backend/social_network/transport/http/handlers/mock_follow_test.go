package handlers

import (
	"context"
	"social_network/internal/follow"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFollowService struct {
	mock.Mock
}

func (m *MockFollowService) Follow(ctx context.Context, follow follow.Follow) error {
	args := m.Called(ctx, follow)
	return args.Error(0)
}

func (m *MockFollowService) Unfollow(ctx context.Context, follow follow.Follow) error {
	args := m.Called(ctx, follow)
	return args.Error(0)
}

func (m *MockFollowService) GetFollowers(ctx context.Context, userID uuid.UUID) ([]follow.Follow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]follow.Follow), args.Error(1)
}

func (m *MockFollowService) GetFollowing(ctx context.Context, userID uuid.UUID) ([]follow.Follow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]follow.Follow), args.Error(1)
}
