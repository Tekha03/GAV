package settings

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		s, err := NewService(repo)
		require.NoError(t, err)
		require.NotNil(t, s)
	})

	t.Run("nil repo", func(t *testing.T) {
		s, err := NewService(nil)
		require.ErrorIs(t, err, ErrRepoNil)
		require.Nil(t, s)
	})
}

func TestService_Get(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	mockSettings := &UserSettings{
		UserID: userID,
		ProfilePrivacy: false,
		ShowLocation: true,
		AllowMessages: true,
	}

	t.Run("success existing settings", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByUserID", ctx, userID).Return(mockSettings, nil).Once()

		s, _ := NewService(repo)
		settings, err := s.Get(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, mockSettings, settings)
	})

	t.Run("nil userID", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)
		settings, err := s.Get(ctx, uuid.Nil)

		require.ErrorIs(t, err, ErrInvalidUserID)
		require.Nil(t, settings)
	})

	t.Run("create default settings if not exists", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByUserID", ctx, userID).Return(nil, ErrSettingsNotFound).Once()
		repo.On("Create", ctx, mock.AnythingOfType("*settings.UserSettings")).Return(nil).Once()

		s, _ := NewService(repo)
		settings, err := s.Get(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, userID, settings.UserID)
		require.False(t, settings.ProfilePrivacy)
		require.True(t, settings.ShowLocation)
		require.True(t, settings.AllowMessages)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	mockSettings := &UserSettings{
		UserID: userID,
		ProfilePrivacy: false,
		ShowLocation: true,
		AllowMessages: true,
	}
	newProfilePrivacy := true
	newShowLocation := false

	t.Run("success update fields", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByUserID", ctx, userID).Return(mockSettings, nil).Once()
		repo.On("Update", ctx, mockSettings).Return(nil).Once()

		s, _ := NewService(repo)
		err := s.Update(ctx, userID, UpdateSettingsInput{
			ProfilePrivacy: &newProfilePrivacy,
			ShowLocation:   &newShowLocation,
		})

		require.NoError(t, err)
		require.Equal(t, newProfilePrivacy, mockSettings.ProfilePrivacy)
		require.Equal(t, newShowLocation, mockSettings.ShowLocation)
	})

	t.Run("invalid userID", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)
		err := s.Update(ctx, uuid.Nil, UpdateSettingsInput{})

		require.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("settings not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByUserID", ctx, userID).Return(nil, ErrSettingsNotFound).Once()

		s, _ := NewService(repo)
		err := s.Update(ctx, userID, UpdateSettingsInput{})

		require.ErrorIs(t, err, ErrSettingsNotFound)
	})
}
