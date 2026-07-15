package auth

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidClaims      = errors.New("invalid claims")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserIdNotFound     = errors.New("user_id not found in token")

	ErrUserServiceNil          = errors.New("auth service: user service is nil")
	ErrTokenServiceNil         = errors.New("auth service: token service is nil")
	ErrJWTSecretNil            = errors.New("auth service: jwt secret is nil")
	ErrHasherNil               = errors.New("auth service: hasher is nil")
	ErrInvalidRefreshToken     = errors.New("auth service: invalid refresh token")
	ErrUserIDNil               = errors.New("claims: user id is nil")
	ErrEmptyRole               = errors.New("claims: role is empty")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)
