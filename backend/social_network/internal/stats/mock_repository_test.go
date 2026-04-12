package stats

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUserStats(ctx context.Context, stats *UserStats) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockRepository) DeleteUserStats(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	args := m.Called(ctx, userID)
	stats := args.Get(0)
	if stats == nil {
		return nil, args.Error(1)
	}
	return stats.(*UserStats), args.Error(1)
}

func (m *MockRepository) IncrementPosts(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) IncrementFollowers(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) IncrementDogs(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) IncrementFollowings(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) DecrementPosts(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) DecrementFollowers(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) DecrementDogs(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) DecrementFollowings(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) CreatePostStats(ctx context.Context, stats *PostStats) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockRepository) GetPostStats(ctx context.Context, postID uuid.UUID) (*PostStats, error) {
	args := m.Called(ctx, postID)
	stats := args.Get(0)
	if stats == nil {
		return nil, args.Error(1)
	}
	return stats.(*PostStats), args.Error(1)
}

func (m *MockRepository) IncrementPostLikes(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockRepository) DecrementPostLikes(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockRepository) IncrementPostComments(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockRepository) DecrementPostComments(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}
