package auth

import "context"

type Service interface {
	Register(ctx context.Context, email, password string) (userID uint, err error)
	Login(ctx context.Context, email, password string) (token string, err error)
}
