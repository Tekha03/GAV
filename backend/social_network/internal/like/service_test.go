package like

import (
	"context"
	"testing"

	"github.com/google/uuid"
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

func TestService_Add(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	postID := uuid.New()

	t.Run("invalid like", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		err := s.Add(ctx, Like{})

		require.ErrorIs(t, err, ErrInvalidLike)
	})

	t.Run("repo error on LikeExists", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		like := Like{UserID: userID, PostID: postID}

		repo.On("LikeExists", ctx, like).Return(false, ErrDBError).Once()

		err := s.Add(ctx, like)

		require.Error(t, err)
	})

	t.Run("already liked", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		like := Like{UserID: userID, PostID: postID}

		repo.On("LikeExists", ctx, like).Return(true, nil).Once()

		err := s.Add(ctx, like)

		require.ErrorIs(t, err, ErrAlreadyLiked)
	})

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		like := Like{UserID: userID, PostID: postID}

		repo.On("LikeExists", ctx, like).Return(false, nil).Once()
		repo.On("Add", ctx, like).Return(nil).Once()

		err := s.Add(ctx, like)

		require.NoError(t, err)
	})
}

func TestService_Remove(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	postID := uuid.New()

	t.Run("invalid like", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		err := s.Remove(ctx, Like{})

		require.ErrorIs(t, err, ErrInvalidLike)
	})

	t.Run("repo error on LikeExists", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		like := Like{UserID: userID, PostID: postID}

		repo.On("LikeExists", ctx, like).Return(false, ErrDBError).Once()

		err := s.Remove(ctx, like)

		require.Error(t, err)
	})

	t.Run("like does not exist", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		like := Like{UserID: userID, PostID: postID}

		repo.On("LikeExists", ctx, like).Return(false, nil).Once()

		err := s.Remove(ctx, like)

		require.ErrorIs(t, err, ErrLikeDoesNotExist)
	})

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		like := Like{UserID: userID, PostID: postID}

		repo.On("LikeExists", ctx, like).Return(true, nil).Once()
		repo.On("Remove", ctx, like).Return(nil).Once()

		err := s.Remove(ctx, like)

		require.NoError(t, err)
	})
}
