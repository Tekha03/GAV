package settings

import apperrors "shared/app_errors"

var (
	ErrSettingsNotFound = apperrors.New(apperrors.SettingsNotFound, "settings not found")
	ErrInvalidUserID    = apperrors.New(apperrors.Validation, "invalid user ID", apperrors.WithDetail("field", "user_id"))

	ErrRepoNil = apperrors.New(apperrors.Internal, "settings service: repo is nil")
)
