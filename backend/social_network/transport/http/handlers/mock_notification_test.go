package handlers

import (
	"context"
	"social_network/internal/notification"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) NotifyLike(ctx context.Context, postOwnerID, likerID, postID uuid.UUID) error {
	args := m.Called(ctx, postOwnerID, likerID, postID)
	return args.Error(0)
}

func (m *MockNotificationService) NotifyComment(ctx context.Context, toUserID, fromUserID, postID uuid.UUID) error {
	args := m.Called(ctx, toUserID, fromUserID, postID)
	return args.Error(0)
}

func (m *MockNotificationService) NotifyFollow(ctx context.Context, followingID, followerID uuid.UUID) error {
	args := m.Called(ctx, followingID, followerID)

	if len(args) == 0 {
		return nil
	}

	return args.Error(0)
}

func (m *MockNotificationService) GetInAppNotifications(ctx context.Context, userID uuid.UUID) ([]*notification.Notification, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*notification.Notification), args.Error(1)
}

func (m *MockNotificationService) MarkInAppNotificationAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	args := m.Called(ctx, userID, notificationID)
	return args.Error(0)
}
