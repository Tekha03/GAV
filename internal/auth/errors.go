package auth

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidClaims = errors.New("invalid claims")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken = errors.New("invalid token")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserIdNotFound = errors.New("user_id not found in token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)
