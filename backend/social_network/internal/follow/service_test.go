package follow

import (
	"context"
	statsPkg "social_network/internal/stats"
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

func TestService_Follow(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)

	s, _ := NewService(repo)

	followerID := uuid.New()
	followingID := uuid.New()

	t.Run("cannot follow yourself", func(t *testing.T) {
		f := Follow{FollowerID: followerID, FollowingID: followerID}

		err := s.Follow(ctx, f)

		require.ErrorIs(t, err, ErrCannotFollowYourself)
	})

	t.Run("already following", func(t *testing.T) {
		f := Follow{FollowerID: followerID, FollowingID: followingID}

		repo.On("FollowerExists", ctx, f).Return(true, nil).Once()

		err := s.Follow(ctx, f)

		require.ErrorIs(t, err, ErrAlreadyFollowing)
	})

	t.Run("repo error on exists", func(t *testing.T) {
		f := Follow{FollowerID: followerID, FollowingID: followingID}

		repo.On("FollowerExists", ctx, f).Return(false, ErrDBError).Once()

		err := s.Follow(ctx, f)

		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		f := Follow{FollowerID: followerID, FollowingID: followingID}

		repo.On("FollowerExists", ctx, f).Return(false, nil).Once()
		repo.On("Follow", ctx, f).Return(nil).Once()

		err := s.Follow(ctx, f)

		require.NoError(t, err)
	})

	t.Run("increments stats in the right direction", func(t *testing.T) {
		repo := new(MockRepository)
		statService := &recordingStatsService{StatsService: statsPkg.NoopService()}

		s, err := NewService(repo, statService)
		require.NoError(t, err)

		f := Follow{FollowerID: followerID, FollowingID: followingID}

		repo.On("FollowerExists", ctx, f).Return(false, nil).Once()
		repo.On("Follow", ctx, f).Return(nil).Once()

		err = s.Follow(ctx, f)

		require.NoError(t, err)
		require.Equal(t, []uuid.UUID{followerID}, statService.incrementedFollowings)
		require.Equal(t, []uuid.UUID{followingID}, statService.incrementedFollowers)
	})
}

func TestService_Unfollow(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)

	s, _ := NewService(repo)

	f := Follow{
		FollowerID:  uuid.New(),
		FollowingID: uuid.New(),
	}

	repo.On("Unfollow", ctx, f).Return(nil).Once()

	err := s.Unfollow(ctx, f)

	require.NoError(t, err)
}

type recordingStatsService struct {
	statsPkg.StatsService
	incrementedFollowers  []uuid.UUID
	incrementedFollowings []uuid.UUID
}

func (s *recordingStatsService) IncrementFollowers(_ context.Context, userID uuid.UUID) error {
	s.incrementedFollowers = append(s.incrementedFollowers, userID)
	return nil
}

func (s *recordingStatsService) IncrementFollowings(_ context.Context, userID uuid.UUID) error {
	s.incrementedFollowings = append(s.incrementedFollowings, userID)
	return nil
}

func TestService_GetFollowers(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)

	s, _ := NewService(repo)

	userID := uuid.New()

	t.Run("invalid user id", func(t *testing.T) {
		res, err := s.GetFollowers(ctx, uuid.Nil)

		require.ErrorIs(t, err, ErrInvalidUserID)
		require.Nil(t, res)
	})

	t.Run("repo error", func(t *testing.T) {
		repo.On("GetFollowers", ctx, userID).Return(nil, ErrDBError).Once()

		res, err := s.GetFollowers(ctx, userID)

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("success", func(t *testing.T) {
		expected := []Follow{
			{FollowerID: uuid.New(), FollowingID: userID},
		}

		repo.On("GetFollowers", ctx, userID).Return(expected, nil).Once()

		res, err := s.GetFollowers(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, expected, res)
	})
}

func TestService_GetFollowing(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)

	s, _ := NewService(repo)

	userID := uuid.New()

	t.Run("invalid user id", func(t *testing.T) {
		res, err := s.GetFollowing(ctx, uuid.Nil)

		require.ErrorIs(t, err, ErrInvalidUserID)
		require.Nil(t, res)
	})

	t.Run("repo error", func(t *testing.T) {
		repo.On("GetFollowing", ctx, userID).Return(nil, ErrDBError).Once()

		res, err := s.GetFollowing(ctx, userID)

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("success", func(t *testing.T) {
		expected := []Follow{
			{FollowerID: userID, FollowingID: uuid.New()},
		}

		repo.On("GetFollowing", ctx, userID).Return(expected, nil).Once()

		res, err := s.GetFollowing(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, expected, res)
	})
}
