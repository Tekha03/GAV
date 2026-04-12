package notification

import (
	"context"
	"social_network/internal/device"
	"social_network/internal/firebase"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	hub := NewHub()
	repo := new(MockNotificationRepo)
	deviceRepo := new(MockDeviceRepo)

	t.Run("success", func(t *testing.T) {
		s, err := NewService(hub, repo, deviceRepo, (*firebase.Client)(nil))
		require.Error(t, err)
		require.Nil(t, s)
	})

	t.Run("nil hub", func(t *testing.T) {
		s, err := NewService(nil, repo, deviceRepo, (*firebase.Client)(nil))
		require.ErrorIs(t, err, ErrFirebaseClientEmpty)
		require.Nil(t, s)
	})
}

func TestService_NotifyLike(t *testing.T) {
	ctx := context.Background()

	hub := NewHub()
	repo := new(MockNotificationRepo)
	deviceRepo := new(MockDeviceRepo)

	fb := &firebase.Client{}

	s, _ := NewService(hub, repo, deviceRepo, fb)

	postOwnerID := uuid.New()
	likerID := uuid.New()
	postID := uuid.New()

	t.Run("repo error", func(t *testing.T) {
		repo.
			On("Create", ctx, mock.Anything).
			Return(ErrNotificationRepoEmpty).
			Once()

		err := s.NotifyLike(ctx, postOwnerID, likerID, postID)

		require.Error(t, err)
	})

	t.Run("success without devices", func(t *testing.T) {
		repo.ExpectedCalls = nil

		repo.
			On("Create", ctx, mock.Anything).
			Return(nil).
			Once()

		deviceRepo.
			On("GetByUser", ctx, postOwnerID).
			Return([]*device.DeviceToken{}, nil).
			Once()

		err := s.NotifyLike(ctx, postOwnerID, likerID, postID)

		require.NoError(t, err)
	})
}

func TestService_NotifyComment(t *testing.T) {
	ctx := context.Background()

	hub := NewHub()
	repo := new(MockNotificationRepo)
	deviceRepo := new(MockDeviceRepo)
	fb := &firebase.Client{}

	s, _ := NewService(hub, repo, deviceRepo, fb)

	postOwnerID := uuid.New()
	commenterID := uuid.New()
	postID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo.
			On("Create", ctx, mock.Anything).
			Return(nil).
			Once()

		deviceRepo.
			On("GetByUser", ctx, postOwnerID).
			Return([]*device.DeviceToken{}, nil).
			Once()

		err := s.NotifyComment(ctx, postOwnerID, commenterID, postID)

		require.NoError(t, err)
	})
}

func TestService_NotifyFollow(t *testing.T) {
	ctx := context.Background()

	hub := NewHub()
	repo := new(MockNotificationRepo)
	deviceRepo := new(MockDeviceRepo)
	fb := &firebase.Client{}

	s, _ := NewService(hub, repo, deviceRepo, fb)

	followingID := uuid.New()
	followerID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo.
			On("Create", ctx, mock.Anything).
			Return(nil).
			Once()

		deviceRepo.
			On("GetByUser", ctx, followingID).
			Return([]*device.DeviceToken{}, nil).
			Once()

		err := s.NotifyFollow(ctx, followingID, followerID)

		require.NoError(t, err)
	})
}

func TestService_GetInAppNotifications(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	repo := new(MockNotificationRepo)

	hub := NewHub()
	fb := &firebase.Client{}

	s, _ := NewService(hub, repo, nil, fb)

	t.Run("repo nil", func(t *testing.T) {
		s2, _ := NewService(hub, nil, nil, fb)

		res, err := s2.GetInAppNotifications(ctx, userID)

		require.ErrorIs(t, err, ErrNotificationRepoEmpty)
		require.Nil(t, res)
	})

	t.Run("success", func(t *testing.T) {
		expected := []*Notification{{UserID: userID}}

		repo.
			On("GetByUserID", ctx, userID).
			Return(expected, nil).
			Once()

		res, err := s.GetInAppNotifications(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, expected, res)
	})
}

func TestService_MarkInAppNotificationAsRead(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	notifID := uuid.New()

	repo := new(MockNotificationRepo)
	hub := NewHub()
	fb := &firebase.Client{}

	s, _ := NewService(hub, repo, nil, fb)

	t.Run("repo nil", func(t *testing.T) {
		s2, _ := NewService(hub, nil, nil, fb)

		err := s2.MarkInAppNotificationAsRead(ctx, userID, notifID)

		require.ErrorIs(t, err, ErrNotificationRepoEmpty)
	})

	t.Run("success", func(t *testing.T) {
		repo.
			On("MarkAsRead", ctx, userID, notifID).
			Return(nil).
			Once()

		err := s.MarkInAppNotificationAsRead(ctx, userID, notifID)

		require.NoError(t, err)
	})
}
