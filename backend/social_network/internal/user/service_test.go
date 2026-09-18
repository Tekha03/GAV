package user

import (
	"context"
	"social_network/internal/dog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type testEnv struct {
	service UserService
	repo    *MockRepository
}

func setup(t *testing.T) *testEnv {
	repo := &MockRepository{}

	service, err := NewService(repo)
	require.NoError(t, err)

	return &testEnv{
		service: service,
		repo:    repo,
	}
}

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

func TestService_Get(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	email := "test@mail.com"
	mockUser := &User{
		ID:       id,
		Email:    email,
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

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("UpdateEmail", ctx, id, email).Return(nil).Once()

		s, _ := NewService(repo)
		err := s.Update(ctx, id, UpdateUserInput{
			Email: &email,
		})

		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("UpdateEmail", ctx, id, email).Return(ErrFail).Once()

		s, _ := NewService(repo)
		err := s.Update(ctx, id, UpdateUserInput{
			Email: &email,
		})

		require.Error(t, err)
	})
}

func TestService_UpdateRejectsMissingOrInvalidEmail(t *testing.T) {
	s, _ := NewService(new(MockRepository))
	id := uuid.New()
	require.ErrorIs(t, s.Update(context.Background(), id, UpdateUserInput{}), ErrUpdateEmpty)
	for _, value := range []string{"", "not-an-email", "Name <name@example.com>"} {
		require.ErrorIs(t, s.Update(context.Background(), id, UpdateUserInput{Email: &value}), ErrEmailInvalid)
	}
}

func TestService_ChangePassword(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	oldHash, err := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.MinCost)
	require.NoError(t, err)
	repo := new(MockRepository)
	repo.On("GetByID", ctx, id).Return(&User{ID: id, Password: string(oldHash)}, nil)
	repo.On("UpdatePassword", ctx, id, string(oldHash), mock.MatchedBy(func(hash string) bool {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte("new-password")) == nil
	})).Return(nil).Once()
	s, _ := NewService(repo)
	require.ErrorIs(t, s.ChangePassword(ctx, id, ChangePasswordInput{CurrentPassword: "wrong", NewPassword: "new-password"}), ErrCurrentPasswordInvalid)
	require.ErrorIs(t, s.ChangePassword(ctx, id, ChangePasswordInput{CurrentPassword: "old-password", NewPassword: "short"}), ErrPasswordInvalid)
	require.NoError(t, s.ChangePassword(ctx, id, ChangePasswordInput{CurrentPassword: "old-password", NewPassword: "new-password"}))
	repo.AssertExpectations(t)
}

func TestService_ChangeRoleRequiresAdmin(t *testing.T) {
	ctx := context.Background()
	actor, target := uuid.New(), uuid.New()
	repo := new(MockRepository)
	repo.On("GetByID", ctx, actor).Return(&User{ID: actor, Role: "user"}, nil).Once()
	repo.On("GetByID", ctx, actor).Return(&User{ID: actor, Role: "admin"}, nil).Twice()
	repo.On("UpdateRole", ctx, target, "admin").Return(nil).Once()
	s, _ := NewService(repo)
	require.ErrorIs(t, s.ChangeRole(ctx, actor, target, ChangeRoleInput{Role: "admin"}), ErrRoleForbidden)
	require.ErrorIs(t, s.ChangeRole(ctx, actor, target, ChangeRoleInput{Role: "owner"}), ErrRoleInvalid)
	require.NoError(t, s.ChangeRole(ctx, actor, target, ChangeRoleInput{Role: "admin"}))
	repo.AssertExpectations(t)
}

func TestService_UpdateLocation_UpsertsActiveWalkSession(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	repo := new(MockRepository)
	repo.On("GetByID", ctx, userID).
		Return(&User{ID: userID, Email: "walk@gav.app", Password: "hash", Role: "user"}, nil).
		Once()
	repo.On("UpsertActiveWalkSession", ctx, mock.MatchedBy(func(session *WalkSession) bool {
		return session != nil &&
			session.UserID == userID &&
			session.Lat == 55.751244 &&
			session.Lon == 37.618423 &&
			session.Visibility == VisibilityEveryone
	})).Return(nil).Once()
	repo.On("Update", ctx, mock.MatchedBy(func(u *User) bool {
		return u != nil &&
			u.ID == userID &&
			u.LocationStatus == Walking &&
			u.Visibility == VisibilityEveryone &&
			u.LocationUpdatedAt == nil
	})).Return(nil).Once()

	s, _ := NewService(repo)
	err := s.UpdateLocation(ctx, userID, UpdateLocationInput{
		Latitude:   55.751244,
		Longitude:  37.618423,
		Status:     Walking,
		Visibility: VisibilityEveryone,
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_UpdateLocation_EndsActiveWalkSessionOnClear(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	repo := new(MockRepository)
	repo.On("GetByID", ctx, userID).
		Return(&User{ID: userID, Email: "walk@gav.app", Password: "hash", Role: "user"}, nil).
		Once()
	repo.On("EndActiveWalkSession", ctx, userID, mock.Anything).Return(nil).Once()
	repo.On("Update", ctx, mock.MatchedBy(func(u *User) bool {
		return u != nil &&
			u.ID == userID &&
			u.Lat == nil &&
			u.Lon == nil &&
			u.LocationUpdatedAt == nil &&
			u.LocationStatus == Inactive &&
			u.Visibility == VisibilityNoOne
	})).Return(nil).Once()

	s, _ := NewService(repo)
	err := s.UpdateLocation(ctx, userID, UpdateLocationInput{
		ClearLocation: true,
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
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

func TestFindDogsNearby_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		env := setup(t)
		otherUserID := uuid.New()

		dogsFromRepo := []*dog.Dog{
			{
				ID:      uuid.New(),
				OwnerID: otherUserID,
			},
		}

		env.repo.
			On("FindWalkingNearby", ctx, 0.0, 0.0, 1000.0, mock.Anything).
			Return(dogsFromRepo, nil)
		env.repo.
			On("GetByID", ctx, otherUserID).
			Return(&User{ID: otherUserID, Visibility: VisibilityEveryone}, nil)

		dogs, err := env.service.FindDogsNearby(ctx, userID, 0.0, 0.0, 1000.0)

		require.NoError(t, err)
		require.Len(t, dogs, 1)
		assert.Equal(t, otherUserID, dogs[0].OwnerID)
	})

	t.Run("exclude own dogs", func(t *testing.T) {
		env := setup(t)
		dogsFromRepo := []*dog.Dog{
			{
				ID:      uuid.New(),
				OwnerID: userID,
			},
		}

		env.repo.
			On("FindWalkingNearby", ctx, 0.0, 0.0, 1000.0, mock.Anything).
			Return(dogsFromRepo, nil)
		env.repo.
			On("GetByID", ctx, userID).
			Return(&User{ID: userID, Visibility: VisibilityEveryone}, nil)

		dogs, err := env.service.FindDogsNearby(ctx, userID, 0.0, 0.0, 1000.0)

		require.NoError(t, err)
		require.Len(t, dogs, 0)
	})

	t.Run("filter invisible", func(t *testing.T) {
		env := setup(t)
		otherUserID := uuid.New()
		dogsFromRepo := []*dog.Dog{
			{
				ID:      uuid.New(),
				OwnerID: otherUserID,
			},
		}

		env.repo.
			On("FindWalkingNearby", ctx, 0.0, 0.0, 1000.0, mock.Anything).
			Return(dogsFromRepo, nil)
		env.repo.
			On("GetByID", ctx, otherUserID).
			Return(&User{ID: otherUserID, Visibility: VisibilityNoOne}, nil)

		dogs, err := env.service.FindDogsNearby(ctx, userID, 0, 0, 1000)

		require.NoError(t, err)
		require.Len(t, dogs, 0)
	})

	t.Run("repo error", func(t *testing.T) {
		env := setup(t)
		env.repo.
			On("FindWalkingNearby", ctx, 0.0, 0.0, 1000.0, mock.Anything).
			Return(nil, assert.AnError)

		dogs, err := env.service.FindDogsNearby(ctx, userID, 0, 0, 1000)

		require.Error(t, err)
		assert.Nil(t, dogs)
	})

	t.Run("invalid coordinates", func(t *testing.T) {
		env := setup(t)

		dogs, err := env.service.FindDogsNearby(ctx, userID, 91, 0, 1000)

		require.ErrorIs(t, err, ErrInvalidLatitude)
		require.Nil(t, dogs)
	})

	t.Run("invalid radius", func(t *testing.T) {
		env := setup(t)

		dogs, err := env.service.FindDogsNearby(ctx, userID, 0, 0, 49)

		require.ErrorIs(t, err, ErrInvalidRadius)
		require.Nil(t, dogs)
	})
}
