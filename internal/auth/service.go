package auth

import "context"

type AuthService interface {
	Register(ctx context.Context, email, password string) (token string, err error)
	Login(ctx context.Context, email, password string) (token string, err error)
}
