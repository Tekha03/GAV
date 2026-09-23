package middleware

import apperrors "shared/app_errors"

var (
	ErrUnauthorized = apperrors.New(apperrors.AuthTokenMissing, "unauthorized")
	ErrForbidden    = apperrors.New(apperrors.AuthForbidden, "insufficient permissions")
	ErrNotOwner     = apperrors.New(apperrors.PostAccessDenied, "not post owner")
	ErrInvalidID    = apperrors.New(apperrors.Validation, "invalid id", apperrors.WithDetail("field", "id"))
)
