package handlers

import (
	"context"
	"social_network/internal/settings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSettingsServcie struct {
	mock.Mock
}

func (m *MockSettingsServcie) Get(ctx context.Context, userID uuid.UUID) (*settings.UserSettings, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*settings.UserSettings), args.Error(1)
}

func (m *MockSettingsServcie) Update(ctx context.Context, userID uuid.UUID, input settings.UpdateSettingsInput) error {
	args := m.Called(ctx, userID, input)
	return args.Error(0)
}
