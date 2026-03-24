package vaccination

import "errors"

var (
	ErrVaccAccessDenied = errors.New("vaccination access denied")
	ErrDogIDEmpty		= errors.New("dogID cannot be empty")

	ErrRepoNil			= errors.New("vaccination service: repo is nil")
)
