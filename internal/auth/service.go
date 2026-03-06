package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*AuthTokens, error)
	Login(ctx context.Context, email, password string) (*AuthTokens, error)
	Me(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
	Logout(ctx context.Context, token string) error
	Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error)
}
