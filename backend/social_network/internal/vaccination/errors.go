package vaccination

import "errors"

var (
	ErrVaccAccessDenied    = errors.New("vaccination access denied")
	ErrDogIDEmpty          = errors.New("dogID cannot be empty")
	ErrDBError             = errors.New("db error")
	ErrRepoNil             = errors.New("vaccination service: repo is nil")
	ErrVaccinationNotFound = errors.New("vaccination not found")
)
