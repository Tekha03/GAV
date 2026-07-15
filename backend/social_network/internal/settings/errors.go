package settings

import "errors"

var (
	ErrSettingsNotFound = errors.New("settings not found")
	ErrInvalidUserID    = errors.New("invalid user ID")

	ErrRepoNil = errors.New("settings service: repo is nil")
)
