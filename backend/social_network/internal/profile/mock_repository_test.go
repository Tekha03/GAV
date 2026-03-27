package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, profile *UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, profileID uuid.UUID) (*UserProfile, error) {
	args := m.Called(ctx, profileID)
	profile := args.Get(0)

	if profile == nil {
		return nil, args.Error(1)
	}

	return profile.(*UserProfile), args.Error(1)
}

func (m *MockRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	args := m.Called(ctx, userID)
	profile := args.Get(0)

	if profile == nil {
		return nil, args.Error(1)
	}

	return profile.(*UserProfile), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, profile *UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, profileID uuid.UUID) error {
	args := m.Called(ctx, profileID)
	return args.Error(0)
}
