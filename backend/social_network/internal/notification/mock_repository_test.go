package notification

import (
	"context"
	"social_network/internal/device"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// ---- MockNotificationRepo ----

type MockNotificationRepo struct {
	mock.Mock
}

func (m *MockNotificationRepo) Create(ctx context.Context, n *Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func (m *MockNotificationRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Notification, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*Notification), args.Error(1)
}

func (m *MockNotificationRepo) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	args := m.Called(ctx, userID, notificationID)
	return args.Error(0)
}

func (m *MockNotificationRepo) DeleteOld(ctx context.Context, userID uuid.UUID, ttl time.Duration) error {
	args := m.Called(ctx, userID, ttl)
	return args.Error(0)
}

// ---- MockDeviceRepo ----

type MockDeviceRepo struct {
	mock.Mock
}

func (m *MockDeviceRepo) Create(ctx context.Context, t *device.DeviceToken) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockDeviceRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]*device.DeviceToken, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*device.DeviceToken), args.Error(1)
}

// ---- MockFirebaseClient ----

type MockFirebaseClient struct {
	mock.Mock
}

func (m *MockFirebaseClient) SendPush(ctx context.Context, token string, title string, body string, data map[string]string) error {
	args := m.Called(ctx, token, title, body, data)
	return args.Error(0)
}
