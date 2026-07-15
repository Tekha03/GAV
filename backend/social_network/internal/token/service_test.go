package token

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
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

func TestService_CreateRefresh(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Create", ctx, mock.AnythingOfType("*token.RefreshToken")).Return(nil).Once()

		s, _ := NewService(repo)
		token, err := s.CreateRefresh(ctx, userID)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Create", ctx, mock.AnythingOfType("*token.RefreshToken")).Return(ErrFail).Once()

		s, _ := NewService(repo)
		token, err := s.CreateRefresh(ctx, userID)
		require.Error(t, err)
		require.Empty(t, token)
	})
}

func TestService_ValidateAndRotate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	refreshPlain := "plain-refresh"
	hash, _ := bcrypt.GenerateFromPassword([]byte(refreshPlain), bcrypt.DefaultCost)

	mockToken := &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: string(hash),
		ExpiresAt: time.Now().Add(time.Hour),
		Revoked:   false,
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByHash", ctx, mock.Anything).Return(mockToken, nil).Once()
		repo.On("Revoke", ctx, mock.Anything).Return(nil).Once()
		repo.On("Create", ctx, mock.AnythingOfType("*token.RefreshToken")).Return(nil).Once()

		s, _ := NewService(repo)
		uid, newToken, err := s.ValidateAndRotate(ctx, refreshPlain)
		require.NoError(t, err)
		require.Equal(t, userID, uid)
		require.NotEmpty(t, newToken)
		repo.AssertExpectations(t)
	})

	t.Run("invalid refresh", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByHash", ctx, mock.Anything).Return(nil, ErrTokenNotFound).Once()

		s, _ := NewService(repo)
		uid, newToken, err := s.ValidateAndRotate(ctx, "bad-refresh")
		require.ErrorIs(t, err, ErrInvalidRefresh)
		require.Equal(t, uuid.Nil, uid)
		require.Empty(t, newToken)
	})
}

func TestService_Revoke(t *testing.T) {
	ctx := context.Background()
	refresh := "refresh-token"

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Revoke", ctx, mock.Anything).Return(nil).Once()
		s, _ := NewService(repo)
		err := s.Revoke(ctx, refresh)
		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Revoke", ctx, mock.Anything).Return(ErrFail).Once()
		s, _ := NewService(repo)
		err := s.Revoke(ctx, refresh)
		require.Error(t, err)
	})
}

func TestService_RevokeAllForUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("RevokeAllForUser", ctx, userID).Return(nil).Once()
		s, _ := NewService(repo)
		err := s.RevokeAllForUser(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("RevokeAllForUser", ctx, userID).Return(ErrFail).Once()
		s, _ := NewService(repo)
		err := s.RevokeAllForUser(ctx, userID)
		require.Error(t, err)
	})
}
