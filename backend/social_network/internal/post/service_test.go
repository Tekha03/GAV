package post

import (
	"context"
	"testing"
	"time"

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
	content := "Hello World"
	image := "image.jpg"

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Create", ctx, mock.AnythingOfType("*post.Post")).Return(nil).Once()

		s, _ := NewService(repo)
		post, err := s.Create(ctx, userID, content, image)

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, post.ID)
		require.Equal(t, userID, post.UserID)
		require.Equal(t, content, post.Content)
		require.Equal(t, image, post.ImageUrl)
	})

	t.Run("empty content", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		post, err := s.Create(ctx, userID, "", image)

		require.ErrorIs(t, err, ErrEmptyContent)
		require.Nil(t, post)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Create", ctx, mock.Anything).Return(ErrRepoNil).Once()

		s, _ := NewService(repo)
		post, err := s.Create(ctx, userID, content, image)

		require.ErrorIs(t, err, ErrRepoNil)
		require.Nil(t, post)
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()
	mockPost := &Post{ID: postID, UserID: userID, Content: "Hi"}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, postID).Return(mockPost, nil).Once()

		s, _ := NewService(repo)
		post, err := s.GetByID(ctx, postID)

		require.NoError(t, err)
		require.Equal(t, mockPost, post)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, postID).Return(nil, nil).Once()

		s, _ := NewService(repo)
		post, err := s.GetByID(ctx, postID)

		require.ErrorIs(t, err, ErrPostNotFound)
		require.Nil(t, post)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, postID).Return(nil, ErrRepoNil).Once()

		s, _ := NewService(repo)
		post, err := s.GetByID(ctx, postID)

		require.ErrorIs(t, err, ErrRepoNil)
		require.Nil(t, post)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	postID := uuid.New()
	mockPost := &Post{ID: postID, UserID: userID, Content: "Hi"}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, postID).Return(mockPost, nil).Once()
		repo.On("Delete", ctx, postID).Return(nil).Once()

		s, _ := NewService(repo)
		err := s.Delete(ctx, userID, postID)

		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, postID).Return(nil, nil).Once()

		s, _ := NewService(repo)
		err := s.Delete(ctx, userID, postID)

		require.ErrorIs(t, err, ErrPostNotFound)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, postID).Return(mockPost, nil).Once()

		s, _ := NewService(repo)
		err := s.Delete(ctx, otherUserID, postID)

		require.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("repo delete error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, postID).Return(mockPost, nil).Once()
		repo.On("Delete", ctx, postID).Return(ErrRepoNil).Once()

		s, _ := NewService(repo)
		err := s.Delete(ctx, userID, postID)

		require.ErrorIs(t, err, ErrRepoNil)
	})
}

func TestService_GetFeed(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()
	posts := []*Post{
		{ID: uuid.New(), CreatedAt: now},
		{ID: uuid.New(), CreatedAt: now.Add(-time.Hour)},
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("ListFeed", ctx, userID, mock.Anything, 2).Return(posts, nil).Once()

		s, _ := NewService(repo)
		res, cursor, err := s.GetFeed(ctx, userID, now, 2)

		require.NoError(t, err)
		require.Equal(t, posts, res)
		require.Equal(t, time.Time{}, cursor) // нет следующей страницы
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("ListFeed", ctx, userID, mock.Anything, 2).Return(nil, ErrRepoNil).Once()

		s, _ := NewService(repo)
		res, cursor, err := s.GetFeed(ctx, userID, now, 2)

		require.ErrorIs(t, err, ErrRepoNil)
		require.Nil(t, res)
		require.Equal(t, time.Time{}, cursor)
	})
}

func TestService_ListByUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	posts := []*Post{
		{ID: uuid.New(), UserID: userID, Content: "Hi"},
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("ListByUser", ctx, userID).Return(posts, nil).Once()

		s, _ := NewService(repo)
		res, err := s.ListByUser(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, posts, res)
	})
}
