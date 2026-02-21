package auth

import (
	"context"

	"gav/internal/user"

	"github.com/google/uuid"
)

type service struct {
	userService user.UserService
	jwtConfig 	 JWTConfig
	hasher		PasswordHasher
}

func NewService(userService user.UserService, jwtConfig JWTConfig, hasher PasswordHasher) AuthService {
	return &service{userService: userService, jwtConfig: jwtConfig, hasher: hasher}
}

func (s *service) Register(ctx context.Context, email, password string) (string, error) {
	_, err := s.userService.GetByEmail(ctx, email)
	if err == nil {
		return "", ErrUserAlreadyExists
	}

	hashedPassword, err := s.hasher.HashPassword(password)
	if err != nil {
		return "", err
	}

	newUser, err := s.userService.Create(ctx, email, hashedPassword)
	if err != nil {
		return "", ErrEmailAlreadyExists
	}

	token, err := GenerateToken(newUser.ID, s.jwtConfig)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return token, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	authorizedUser, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !s.hasher.CheckPassword(authorizedUser.Password, password) {
		return "", ErrInvalidCredentials
	}

	token, err := GenerateToken(authorizedUser.ID, s.jwtConfig)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) Me(ctx context.Context, userID uuid.UUID) (*UserInfo, error) {
	user, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInfo := UserInfo{ID: userID, Email: user.Email}
	return &userInfo, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	return nil
}
