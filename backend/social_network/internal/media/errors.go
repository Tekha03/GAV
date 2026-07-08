package media

import "errors"

var (
	ErrStorageNil      = errors.New("media service: storage is nil")
	ErrFileTooLarge    = errors.New("file too large (max 5MB)")
	ErrInvalidFileType = errors.New("inly jpg, png, webp allowed")
	ErrInvalidURL      = errors.New("invalid url")
	ErrEmptyURL        = errors.New("empty url")
)
