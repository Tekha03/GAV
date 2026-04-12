package settings

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, settings *UserSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*UserSettings, error) {
	args := m.Called(ctx, userID)
	settings := args.Get(0)

	if settings == nil {
		return nil, args.Error(1)
	}

	return settings.(*UserSettings), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, settings *UserSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}
