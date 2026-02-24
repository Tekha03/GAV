package auth

import "errors"

var (
	ErrEmailAlreadyExists 		= errors.New("email already exists")
	ErrInvalidClaims 			= errors.New("invalid claims")
	ErrInvalidCredentials 		= errors.New("invalid credentials")
	ErrInvalidToken 			= errors.New("invalid token")
	ErrUserAlreadyExists 		= errors.New("user already exists")
	ErrUserIdNotFound 			= errors.New("user_id not found in token")
	ErrUnexpectedSigningMethod 	= errors.New("unexpected signing method")

	ErrUserServiceNil 			= errors.New("auth service: userService is nil")
	ErrJWTSecretNil 			= errors.New("auth service: jwt secret is nil")
	ErrHasherNil 				= errors.New("auth service: hasher is nil")
)
