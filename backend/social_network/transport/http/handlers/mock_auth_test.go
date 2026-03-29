package handlers

import (
	"context"
	"social_network/internal/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, email, password string) (*auth.AuthTokens, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*auth.AuthTokens), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (*auth.AuthTokens, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*auth.AuthTokens), args.Error(1)
}

func (m *MockAuthService) Me(ctx context.Context, userID uuid.UUID) (*auth.UserInfo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*auth.UserInfo), args.Error(1)
}

func (m *MockAuthService) Refresh(ctx context.Context, token string) (*auth.AuthTokens, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*auth.AuthTokens), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
