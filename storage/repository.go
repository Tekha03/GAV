package storage

import (
	"context"

	"gav/internal/user"
)

type Repository interface {
	Create(ctx context.Context, user *user.User) error
	GetByID(ctx context.Context, id uint) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id uint) error
}
