package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Me(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
	Logout(ctx context.Context, token string) error
}
