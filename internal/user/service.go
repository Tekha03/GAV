package user

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
    Create(ctx context.Context, email, passwordHash string) (*User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, id uuid.UUID, input UpdateuserInput) error
    Delete(ctx context.Context, id uuid.UUID) error
}
