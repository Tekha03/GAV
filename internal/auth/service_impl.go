package auth

import (
	"context"
	"time"

	"gav/storage"
	"gav/user"
)

type AuthService struct {
	users storage.Repository
}

func NewAuthService(users storage.Repository) *AuthService {
	return &AuthService{users: users}
}

func (as *AuthService) Register(ctx context.Context, email, password string) (string, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	newUser := user.NewUser(
		0,
		nil,
		nil,
		&user.UserSettings{
			Email: email,
			PasswordHash: hashedPassword,
			CreatedAt: time.Now(),
		},
		nil,
	)

	if err := as.users.Create(ctx, newUser); err != nil {
		return "", ErrEmailAlreadyExists
	}

	token, err := GenerateToken(int(newUser.ID))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return token, nil
}

func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	authorizedUser, err := as.users.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !CheckPassword(authorizedUser.Settings.PasswordHash, password) {
		return "", ErrInvalidCredentials
	}

	token, err := GenerateToken(int(authorizedUser.ID))
	if err != nil {
		return "", err
	}

	return token, nil
}
