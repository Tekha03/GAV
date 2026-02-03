package user

import "context"

type UserService interface {
    GetByID(ctx context.Context, id uint) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, id uint, input UpdateuserInput)
    Delete(ctx context.Context, id uint) error
}
