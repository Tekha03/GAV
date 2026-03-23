package profile

import "errors"

var (
	ErrProfileAlreadyExists = errors.New("profile already exists")
	ErrProfileNotFound		= errors.New("profile not found")
	ErrInvalidUserID	   = errors.New("invalid user ID")
	ErrInvalidProfileID	  	= errors.New("invalid profile ID")

	ErrRepoNil			   = errors.New("profile service: repo is nil")
)
