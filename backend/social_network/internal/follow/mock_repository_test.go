package follow

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Follow(ctx context.Context, f Follow) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockRepository) Unfollow(ctx context.Context, f Follow) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockRepository) FollowerExists(ctx context.Context, f Follow) (bool, error) {
	args := m.Called(ctx, f)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]Follow, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Follow), args.Error(1)
}

func (m *MockRepository) GetFollowing(ctx context.Context, userID uuid.UUID) ([]Follow, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Follow), args.Error(1)
}
