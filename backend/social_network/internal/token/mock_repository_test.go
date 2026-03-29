package token

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, t *RefreshToken) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockRepository) GetByHash(ctx context.Context, hash string) (*RefreshToken, error) {
	args := m.Called(ctx, hash)
	t := args.Get(0)
	if t == nil {
		return nil, args.Error(1)
	}
	return t.(*RefreshToken), args.Error(1)
}

func (m *MockRepository) Revoke(ctx context.Context, hash string) error {
	args := m.Called(ctx, hash)
	return args.Error(0)
}

func (m *MockRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
