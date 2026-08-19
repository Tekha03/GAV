package profile

import apperrors "shared/app_errors"

var (
	ErrProfileAlreadyExists = apperrors.New(apperrors.ProfileAlreadyExists, "profile already exists")
	ErrProfileNotFound      = apperrors.New(apperrors.ProfileNotFound, "profile not found")
	ErrInvalidUserID        = apperrors.New(apperrors.Validation, "invalid user ID", apperrors.WithDetail("field", "user_id"))
	ErrInvalidProfileID     = apperrors.New(apperrors.Validation, "invalid profile ID", apperrors.WithDetail("field", "profile_id"))

	ErrRepoNil = apperrors.New(apperrors.Internal, "profile service: repo is nil")
)
