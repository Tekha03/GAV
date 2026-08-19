package user

import apperrors "shared/app_errors"

var (
	ErrEmailEmpty        = apperrors.New(apperrors.UserEmailRequired, "email is required")
	ErrPasswordHashEmpty = apperrors.New(apperrors.UserPasswordRequired, "password is required")
	ErrRepoNil           = apperrors.New(apperrors.Internal, "user service: repo is nil")
	ErrUserNotFound      = apperrors.New(apperrors.UserNotFound, "user not found")
	ErrFail              = apperrors.New(apperrors.Internal, "fail")
)
