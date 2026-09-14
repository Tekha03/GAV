package user

import apperrors "shared/app_errors"

var (
	ErrEmailEmpty                = apperrors.New(apperrors.UserEmailRequired, "email is required")
	ErrPasswordHashEmpty         = apperrors.New(apperrors.UserPasswordRequired, "password is required")
	ErrRepoNil                   = apperrors.New(apperrors.Internal, "user service: repo is nil")
	ErrUserNotFound              = apperrors.New(apperrors.UserNotFound, "user not found")
	ErrFail                      = apperrors.New(apperrors.Internal, "fail")
	ErrInvalidLatitude           = apperrors.New(apperrors.Validation, "latitude must be between -90 and 90")
	ErrInvalidLongitude          = apperrors.New(apperrors.Validation, "longitude must be between -180 and 180")
	ErrInvalidRadius             = apperrors.New(apperrors.Validation, "radius must be between 50 and 10000 meters")
	ErrInvalidLocationStatus     = apperrors.New(apperrors.Validation, "invalid location status")
	ErrInvalidLocationVisibility = apperrors.New(apperrors.Validation, "invalid location visibility")
)
