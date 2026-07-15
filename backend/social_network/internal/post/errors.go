package post

import "errors"

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden")
	ErrEmptyContent = errors.New("empty content")

	ErrRepoNil = errors.New("post service: repo is nil")
)
