package settings

import appErrors "social_network/internal/errors"

var (
	ErrSettingsNotFound = appErrors.New(appErrors.CodeNotFound, "settings not found")
	ErrInvalidUserID    = appErrors.New(appErrors.CodeValidation, "invalid user ID")

	ErrRepoNil = appErrors.New(appErrors.CodeInternal, "settings service: repo is nil")
)
