package stats

import (
	"context"
	"errors"
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

func TestService_UserStats(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	mockStats := &UserStats{
		UserID: userID,
		PostCount: 5,
		Followers: 10,
		Followings: 3,
		DogsCount: 2,
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetUserStats", ctx, userID).Return(mockStats, nil).Once()

		s, _ := NewService(repo)
		stats, err := s.UserStats(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, mockStats, stats)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetUserStats", ctx, userID).Return(nil, ErrStatsNotFound).Once()

		s, _ := NewService(repo)
		stats, err := s.UserStats(ctx, userID)
		require.Error(t, err)
		require.Nil(t, stats)
	})
}

func TestService_ProfileStats(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	mockStats := &UserStats{
		UserID: userID,
		PostCount: 3,
		Followers: 5,
		Followings: 2,
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetUserStats", ctx, userID).Return(mockStats, nil).Once()

		s, _ := NewService(repo)
		profileStats, err := s.ProfileStats(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, mockStats.UserID, profileStats.UserID)
		require.Equal(t, mockStats.PostCount, profileStats.PostCount)
		require.Equal(t, mockStats.Followers, profileStats.Followers)
		require.Equal(t, mockStats.Followings, profileStats.Followings)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetUserStats", ctx, userID).Return(nil, ErrStatsNotFound).Once()

		s, _ := NewService(repo)
		profileStats, err := s.ProfileStats(ctx, userID)
		require.ErrorIs(t, err, ErrStatsNotFound)
		require.Nil(t, profileStats)
	})
}

func TestService_IncrementDecrement(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	postID := uuid.New()

	repo := new(MockRepository)
	s, _ := NewService(repo)

	methods := []struct{
		name string
		call func() error
		mockCall func()
	}{
		{
			"IncrementPosts",
			func() error {
				return s.IncrementPosts(ctx, userID)
			},
			func() {
				repo.On("IncrementPosts", ctx, userID).Return(nil).Once()
			},
		},
		{
			"IncrementFollowers",
			func() error {
				return s.IncrementFollowers(ctx, userID)
			},
			func() {
				repo.On("IncrementFollowers", ctx, userID).Return(nil).Once()
			},
		},
		{
			"IncrementDogs",
			func() error {
				return s.IncrementDogs(ctx, userID)
			},
			func() {
				repo.On("IncrementDogs", ctx, userID).Return(nil).Once()
			},
		},
		{
			"IncrementFollowings",
			func() error {
				return s.IncrementFollowings(ctx, userID)
			},
			func() {
				repo.On("IncrementFollowings", ctx, userID).Return(nil).Once()
			},
		},
		{
			"DecrementPosts",
			func() error {
				return s.DecrementPosts(ctx, userID)
			},
			func() {
				repo.On("DecrementPosts", ctx, userID).Return(nil).Once()
			},
		},
		{
			"DecrementFollowers",
			func() error {
				return s.DecrementFollowers(ctx, userID)
			},
			func() {
				repo.On("DecrementFollowers", ctx, userID).Return(nil).Once()
			},
		},
		{
			"DecrementDogs",
			func() error {
				return s.DecrementDogs(ctx, userID)
			},
			func() {
				repo.On("DecrementDogs", ctx, userID).Return(nil).Once()
			},
		},
		{
			"DecrementFollowings",
			func() error {
				return s.DecrementFollowings(ctx, userID)
			},
			func() {
				repo.On("DecrementFollowings", ctx, userID).Return(nil).Once()
			},
		},
		{
			"IncrementPostLikes",
			func() error {
				return s.IncrementPostLikes(ctx, postID)
			},
			func() {
				repo.On("IncrementPostLikes", ctx, postID).Return(nil).Once()
			},
		},
		{
			"DecrementPostLikes",
			func() error {
				return s.DecrementPostLikes(ctx, postID)
			},
			func() {
				repo.On("DecrementPostLikes", ctx, postID).Return(nil).Once()
			},
		},
		{
			"IncrementPostComments",
			func() error {
				return s.IncrementPostComments(ctx, postID)
			},
			func() {
				repo.On("IncrementPostComments", ctx, postID).Return(nil).Once()
			},
		},
		{
			"DecrementPostComments",
			func() error {
				return s.DecrementPostComments(ctx, postID)
			},
			func() {
				repo.On("DecrementPostComments", ctx, postID).Return(nil).Once()
			},
		},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			m.mockCall()
			err := m.call()
			require.NoError(t, err)
		})
	}
}

func TestService_PostStats(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	mockPostStats := &PostStats{
		PostID: postID,
		LikesCount: 10,
		CommentsCount: 5,
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetPostStats", ctx, postID).Return(mockPostStats, nil).Once()
		s, _ := NewService(repo)

		stats, err := s.PostStats(ctx, postID)
		require.NoError(t, err)
		require.Equal(t, mockPostStats, stats)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		repo.On("GetPostStats", ctx, postID).Return(nil, errors.New("not found")).Once()
		s, _ := NewService(repo)

		stats, err := s.PostStats(ctx, postID)
		require.Error(t, err)
		require.Nil(t, stats)
	})
}
