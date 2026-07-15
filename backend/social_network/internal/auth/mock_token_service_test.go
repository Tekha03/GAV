package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) CreateRefresh(ctx context.Context, id uuid.UUID) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateAndRotate(ctx context.Context, token string) (uuid.UUID, string, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uuid.UUID), args.String(1), args.Error(2)
}

func (m *MockTokenService) Revoke(ctx context.Context, hash string) error {
	args := m.Called(ctx, hash)
	return args.Error(0)
}

func (m *MockTokenService) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
