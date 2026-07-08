package profile

import appErrors "social_network/internal/errors"

var (
	ErrProfileAlreadyExists = appErrors.New(appErrors.CodeConflict, "profile already exists")
	ErrProfileNotFound      = appErrors.New(appErrors.CodeNotFound, "profile not found")
	ErrInvalidUserID        = appErrors.New(appErrors.CodeValidation, "invalid user ID")
	ErrInvalidProfileID     = appErrors.New(appErrors.CodeValidation, "invalid profile ID")

	ErrRepoNil = appErrors.New(appErrors.CodeInternal, "profile service: repo is nil")
)
