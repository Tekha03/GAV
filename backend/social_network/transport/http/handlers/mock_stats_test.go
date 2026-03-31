package handlers

import (
	"context"
	"social_network/internal/stats"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockStatsService struct {
	mock.Mock
}

func (m *MockStatsService) UserStats(ctx context.Context, userID uuid.UUID) (*stats.UserStats, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*stats.UserStats), args.Error(1)
}

func (m *MockStatsService) ProfileStats(ctx context.Context, userID uuid.UUID) (*stats.ProfileStats, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*stats.ProfileStats), args.Error(1)
}

func (m *MockStatsService) IncrementPosts(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStatsService) IncrementFollowers(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStatsService) IncrementDogs(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStatsService) IncrementFollowings(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStatsService) DecrementPosts(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStatsService) DecrementFollowers(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStatsService) DecrementDogs(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStatsService) DecrementFollowings(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStatsService) PostStats(ctx context.Context, postID uuid.UUID) (*stats.PostStats, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(*stats.PostStats), args.Error(1)
}

func (m *MockStatsService) IncrementPostLikes(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockStatsService) IncrementPostComments(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockStatsService) DecrementPostLikes(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockStatsService) DecrementPostComments(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}
