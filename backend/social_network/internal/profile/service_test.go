package profile

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

func TestService_Create(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	input := CreateProfileInput{
		Name:            "John",
		Surname:         "Doe",
		Username:        "johndoe",
		ProfilePhotoUrl: "url",
		Bio:             "bio",
		Address:         "addr",
		BirthDate:       "2000-01-01",
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Create", ctx, mock.AnythingOfType("*profile.UserProfile")).Return(nil).Once()

		s, _ := NewService(repo)
		profile, err := s.Create(ctx, userID, input)

		require.NoError(t, err)
		require.Equal(t, userID, profile.UserID)
		require.Equal(t, input.Name, profile.Name)
		require.Equal(t, input.Username, profile.Username)
	})

	t.Run("invalid userID", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)
		profile, err := s.Create(ctx, uuid.Nil, input)

		require.ErrorIs(t, err, ErrInvalidUserID)
		require.Nil(t, profile)
	})

	t.Run("profile already exists", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Create", ctx, mock.Anything).Return(ErrProfileAlreadyExists).Once()

		s, _ := NewService(repo)
		profile, err := s.Create(ctx, userID, input)

		require.ErrorIs(t, err, ErrProfileAlreadyExists)
		require.Nil(t, profile)
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	profileID := uuid.New()
	mockProfile := &UserProfile{UserID: profileID, Name: "John"}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, profileID).Return(mockProfile, nil).Once()

		s, _ := NewService(repo)
		profile, err := s.GetByID(ctx, profileID)

		require.NoError(t, err)
		require.Equal(t, mockProfile, profile)
	})

	t.Run("invalid profileID", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)
		profile, err := s.GetByID(ctx, uuid.Nil)

		require.ErrorIs(t, err, ErrInvalidProfileID)
		require.Nil(t, profile)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, profileID).Return(nil, ErrProfileNotFound).Once()

		s, _ := NewService(repo)
		profile, err := s.GetByID(ctx, profileID)

		require.ErrorIs(t, err, ErrProfileNotFound)
		require.Nil(t, profile)
	})
}

func TestService_GetByUserID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	mockProfile := &UserProfile{UserID: userID, Name: "John"}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByUserID", ctx, userID).Return(mockProfile, nil).Once()

		s, _ := NewService(repo)
		profile, err := s.GetByUserID(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, mockProfile, profile)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByUserID", ctx, userID).Return(nil, ErrProfileNotFound).Once()

		s, _ := NewService(repo)
		profile, err := s.GetByUserID(ctx, userID)

		require.Error(t, err)
		require.Nil(t, profile)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	profileID := uuid.New()
	mockProfile := &UserProfile{UserID: profileID, Name: "John"}
	newName := "Jane"

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, profileID).Return(mockProfile, nil).Once()
		repo.On("Update", ctx, mockProfile).Return(nil).Once()

		s, _ := NewService(repo)
		err := s.Update(ctx, profileID, UpdateProfileInput{Name: &newName})

		require.NoError(t, err)
		require.Equal(t, newName, mockProfile.Name)
	})

	t.Run("invalid profileID", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)
		err := s.Update(ctx, uuid.Nil, UpdateProfileInput{Name: &newName})

		require.ErrorIs(t, err, ErrInvalidProfileID)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, profileID).Return(nil, ErrProfileNotFound).Once()

		s, _ := NewService(repo)
		err := s.Update(ctx, profileID, UpdateProfileInput{Name: &newName})

		require.ErrorIs(t, err, ErrProfileNotFound)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	profileID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Delete", ctx, profileID).Return(nil).Once()

		s, _ := NewService(repo)
		err := s.Delete(ctx, profileID)

		require.NoError(t, err)
	})

	t.Run("invalid profileID", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)
		err := s.Delete(ctx, uuid.Nil)

		require.ErrorIs(t, err, ErrInvalidProfileID)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Delete", ctx, profileID).Return(ErrProfileNotFound).Once()

		s, _ := NewService(repo)
		err := s.Delete(ctx, profileID)

		require.ErrorIs(t, err, ErrProfileNotFound)
	})
}
