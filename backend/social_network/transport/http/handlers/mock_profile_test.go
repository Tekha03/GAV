package handlers

import (
	"context"
	"social_network/internal/profile"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) Create(ctx context.Context, userID uuid.UUID, input profile.CreateProfileInput) (*profile.UserProfile, error) {
	args := m.Called(ctx, userID, input)
	return args.Get(0).(*profile.UserProfile), args.Error(1)
}

func (m *MockProfileService) GetByID(ctx context.Context, profileID uuid.UUID) (*profile.UserProfile, error) {
	args := m.Called(ctx, profileID)
	return args.Get(0).(*profile.UserProfile), args.Error(1)
}

func (m *MockProfileService) GetByUserID(ctx context.Context, userID uuid.UUID) (*profile.UserProfile, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*profile.UserProfile), args.Error(1)
}

func (m *MockProfileService) Update(ctx context.Context, profileID uuid.UUID, input profile.UpdateProfileInput) error {
	args := m.Called(ctx, profileID, input)
	return args.Error(0)
}

func (m *MockProfileService) Delete(ctx context.Context, profileID uuid.UUID) error {
	args := m.Called(ctx, profileID)
	return args.Error(0)
}
