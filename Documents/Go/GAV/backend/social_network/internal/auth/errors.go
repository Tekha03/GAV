package auth

import (
	"errors"

	appErrors "social_network/internal/errors"
)

var (
	ErrEmailAlreadyExists  = appErrors.New(appErrors.CodeConflict, "email already exists")
	ErrInvalidCredentials  = appErrors.New(appErrors.CodeUnauthorized, "invalid email or password")
	ErrUserAlreadyExists   = appErrors.New(appErrors.CodeConflict, "user already exists")
	ErrInvalidRefreshToken = appErrors.New(appErrors.CodeUnauthorized, "invalid refresh token")
	ErrInvalidToken        = appErrors.New(appErrors.CodeUnauthorized, "invalid token")

	ErrInvalidClaims           = errors.New("invalid claims")
	ErrUserIdNotFound          = errors.New("user_id not found in token")
	ErrUserServiceNil          = errors.New("auth service: user service is nil")
	ErrTokenServiceNil         = errors.New("auth service: token service is nil")
	ErrJWTSecretNil            = errors.New("auth service: jwt secret is nil")
	ErrHasherNil               = errors.New("auth service: hasher is nil")
	ErrUserIDNil               = errors.New("claims: user id is nil")
	ErrEmptyRole               = errors.New("claims: role is empty")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)
