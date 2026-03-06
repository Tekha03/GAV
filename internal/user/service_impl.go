package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type service struct {
	repo Repository
}

func NewService(repo Repository) (UserService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo}, nil
}

func (s *service) Create(ctx context.Context, email, passwordHash string) (*User, error) {
	user, err := NewUser(email, passwordHash)
	if err != nil {
		return nil, err
	}

	if err = s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) error {
	user := &User{
		ID: id,
		Email: *input.Email,
		Password: *input.Password,
		Role: *input.Role,
	}

	return s.repo.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
