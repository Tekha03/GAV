package like

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Add(ctx context.Context, like Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *MockRepository) Remove(ctx context.Context, like Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *MockRepository) LikeExists(ctx context.Context, like Like) (bool, error) {
	args := m.Called(ctx, like)
	return args.Bool(0), args.Error(1)
}
