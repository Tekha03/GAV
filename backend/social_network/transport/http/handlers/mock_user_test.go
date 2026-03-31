package handlers

import (
	"context"
	"social_network/internal/dog"
	"social_network/internal/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, email, passwordHash string) (*user.User, error) {
	args := m.Called(ctx, email, passwordHash)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, id uuid.UUID, input user.UpdateUserInput) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) FindDogsNearby(ctx context.Context, id uuid.UUID, centerLat, centerLon float64, radiusMeters float64) ([]*dog.Dog, error) {
	args := m.Called(ctx, id, centerLat, centerLon, radiusMeters)
	return args.Get(0).([]*dog.Dog), args.Error(1)
}
