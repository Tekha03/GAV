package user

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

func TestNewUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		u, err := NewUser("test@mail.com", "hash")
		require.NoError(t, err)
		require.Equal(t, "test@mail.com", u.Email)
		require.Equal(t, "hash", u.Password)
	})

	t.Run("empty email", func(t *testing.T) {
		u, err := NewUser("", "hash")
		require.ErrorIs(t, err, ErrEmailEmpty)
		require.Nil(t, u)
	})

	t.Run("empty password", func(t *testing.T) {
		u, err := NewUser("test@mail.com", "")
		require.ErrorIs(t, err, ErrPasswordHashEmpty)
		require.Nil(t, u)
	})
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()
	email := "test@mail.com"
	password := "hash"

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil).Once()

		s, _ := NewService(repo)
		user, err := s.Create(ctx, email, password)

		require.NoError(t, err)
		require.Equal(t, email, user.Email)
		require.Equal(t, password, user.Password)
	})

	t.Run("invalid input", func(t *testing.T) {
		repo := new(MockRepository)
		s, _ := NewService(repo)

		user, err := s.Create(ctx, "", password)
		require.Error(t, err)
		require.Nil(t, user)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Create", ctx, mock.Anything).Return(ErrFail).Once()

		s, _ := NewService(repo)
		user, err := s.Create(ctx, email, password)

		require.Error(t, err)
		require.Nil(t, user)
	})
}

// --- GetByID / GetByEmail

func TestService_Get(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	email := "test@mail.com"
	mockUser := &User{
		ID: id,
		Email: email,
		Password: "hash",
	}

	t.Run("GetByID success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, id).Return(mockUser, nil).Once()

		s, _ := NewService(repo)
		u, err := s.GetByID(ctx, id)

		require.NoError(t, err)
		require.Equal(t, mockUser, u)
	})

	t.Run("GetByEmail success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByEmail", ctx, email).Return(mockUser, nil).Once()

		s, _ := NewService(repo)
		u, err := s.GetByEmail(ctx, email)

		require.NoError(t, err)
		require.Equal(t, mockUser, u)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetByID", ctx, id).Return(nil, ErrUserNotFound).Once()

		s, _ := NewService(repo)
		u, err := s.GetByID(ctx, id)

		require.Error(t, err)
		require.Nil(t, u)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	email := "new@mail.com"
	password := "newhash"
	role := "admin"

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil).Once()

		s, _ := NewService(repo)
		err := s.Update(ctx, id, UpdateUserInput{
			Email: &email,
			Password: &password,
			Role: &role,
		})

		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Update", ctx, mock.Anything).Return(ErrFail).Once()

		s, _ := NewService(repo)
		err := s.Update(ctx, id, UpdateUserInput{
			Email: &email,
			Password: &password,
			Role: &role,
		})

		require.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Delete", ctx, id).Return(nil).Once()

		s, _ := NewService(repo)
		err := s.Delete(ctx, id)

		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("Delete", ctx, id).Return(ErrFail).Once()

		s, _ := NewService(repo)
		err := s.Delete(ctx, id)

		require.Error(t, err)
	})
}
