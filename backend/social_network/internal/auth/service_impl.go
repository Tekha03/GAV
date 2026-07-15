package auth

import (
	"context"
	"social_network/internal/token"
	"social_network/internal/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userService  user.UserService
	tokenService token.TokenService
	jwtConfig    JWTConfig
	hasher       *PasswordHasher
}

func NewService(userService user.UserService, tokenService token.TokenService, jwtConfig JWTConfig, hasher *PasswordHasher) (AuthService, error) {
	if userService == nil {
		return nil, ErrUserServiceNil
	}
	if tokenService == nil {
		return nil, ErrTokenServiceNil
	}
	if jwtConfig.Secret == nil {
		return nil, ErrJWTSecretNil
	}
	if hasher == nil {
		return nil, ErrHasherNil
	}

	return &service{userService: userService, tokenService: tokenService, jwtConfig: jwtConfig, hasher: hasher}, nil
}

func (s *service) Register(ctx context.Context, email, password string) (*AuthTokens, error) {
	_, err := s.userService.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := s.hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser, err := s.userService.Create(ctx, email, hashedPassword)
	if err != nil {
		return nil, ErrEmailAlreadyExists
	}

	accessToken, err := GenerateAccessToken(newUser.ID, newUser.Role, s.jwtConfig)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	refreshToken, err := s.tokenService.CreateRefresh(ctx, newUser.ID)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Login(ctx context.Context, email, password string) (*AuthTokens, error) {
	authorizedUser, err := s.userService.GetByEmail(ctx, email)
	if err != nil || !s.hasher.CheckPassword(authorizedUser.Password, password) {
		return nil, ErrInvalidCredentials
	}

	access, err := GenerateAccessToken(authorizedUser.ID, authorizedUser.Role, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	refreshStr, err := s.tokenService.CreateRefresh(ctx, authorizedUser.ID)

	return &AuthTokens{
		AccessToken:  access,
		RefreshToken: refreshStr,
	}, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	refrTokenID, _, err := s.tokenService.ValidateAndRotate(ctx, refreshToken)

	user, err := s.userService.GetByID(ctx, refrTokenID)
	if err != nil {
		return nil, err
	}

	newAccess, err := GenerateAccessToken(user.ID, user.Role, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	newRefreshStr, err := s.tokenService.CreateRefresh(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{AccessToken: newAccess, RefreshToken: newRefreshStr}, nil

}

func (s *service) Me(ctx context.Context, userID uuid.UUID) (*UserInfo, error) {
	user, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInfo := UserInfo{ID: userID, Email: user.Email}
	return &userInfo, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), 12)
	return s.tokenService.Revoke(ctx, string(hash))
}
