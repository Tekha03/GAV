package auth

import (
	"context"

	"gav/internal/user"
	"gav/storage"
)

type Service struct {
	users storage.Repository
}

func NewService(users storage.Repository) *Service {
	return &Service{users: users}
}

func (as *Service) Register(ctx context.Context, email, password string) (string, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	newUser := user.NewUser(email, hashedPassword)

	if err := as.users.Create(ctx, newUser); err != nil {
		return "", ErrEmailAlreadyExists
	}

	token, err := GenerateToken(int(newUser.ID))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return token, nil
}

func (as *Service) Login(ctx context.Context, email, password string) (string, error) {
	authorizedUser, err := as.users.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !CheckPassword(authorizedUser.Password, password) {
		return "", ErrInvalidCredentials
	}

	token, err := GenerateToken(int(authorizedUser.ID))
	if err != nil {
		return "", err
	}

	return token, nil
}
