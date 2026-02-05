package user

import (
	"context"
)

type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetByID(ctx context.Context, id uint) (*User, error)
    Update(ctx context.Context, user User) error
    Delete(ctx context.Context, id uint) error
}
