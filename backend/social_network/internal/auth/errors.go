package auth

import apperrors "shared/app_errors"

var (
	ErrEmailAlreadyExists = apperrors.New(apperrors.AuthEmailExists, "email already exists")
	ErrInvalidClaims      = apperrors.New(apperrors.AuthTokenInvalid, "invalid claims")
	ErrInvalidCredentials = apperrors.New(apperrors.AuthCredentialsInvalid, "invalid credentials")
	ErrUserAlreadyExists  = apperrors.New(apperrors.AuthEmailExists, "user already exists")
	ErrUserIdNotFound     = apperrors.New(apperrors.AuthTokenInvalid, "user_id not found in token")

	ErrUserServiceNil          = apperrors.New(apperrors.Internal, "auth service: user service is nil")
	ErrTokenServiceNil         = apperrors.New(apperrors.Internal, "auth service: token service is nil")
	ErrJWTSecretNil            = apperrors.New(apperrors.Internal, "auth service: jwt secret is nil")
	ErrHasherNil               = apperrors.New(apperrors.Internal, "auth service: hasher is nil")
	ErrInvalidRefreshToken     = apperrors.New(apperrors.AuthRefreshInvalid, "invalid refresh token")
	ErrUserIDNil               = apperrors.New(apperrors.AuthTokenInvalid, "claims: user id is nil")
	ErrEmptyRole               = apperrors.New(apperrors.AuthTokenInvalid, "claims: role is empty")
	ErrUnexpectedSigningMethod = apperrors.New(apperrors.AuthTokenInvalid, "unexpected signing method")
	ErrInvalidToken            = apperrors.New(apperrors.AuthTokenInvalid, "invalid token")
)
