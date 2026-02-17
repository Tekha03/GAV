package auth

import (
	"context"

	"gav/internal/user"
	"gav/storage"
)

type service struct {
	users storage.Repository
	jwtConfig JWTConfig
}

func NewService(users storage.Repository, jwtConfig JWTConfig) AuthService {
	return &service{users: users, jwtConfig: jwtConfig}
}

func (s *service) Register(ctx context.Context, email, password string) (string, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	newUser := user.NewUser(email, hashedPassword)

	if err := s.users.Create(ctx, newUser); err != nil {
		return "", ErrEmailAlreadyExists
	}

	token, err := GenerateToken(newUser.ID, s.jwtConfig)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return token, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	authorizedUser, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !CheckPassword(authorizedUser.Password, password) {
		return "", ErrInvalidCredentials
	}

	token, err := GenerateToken(authorizedUser.ID, s.jwtConfig)
	if err != nil {
		return "", err
	}

	return token, nil
}
