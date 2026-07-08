package auth

import (
	"context"

	"social_network/internal/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, email, password string) (*user.User, error) {
	args := m.Called(ctx, email, password)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

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
