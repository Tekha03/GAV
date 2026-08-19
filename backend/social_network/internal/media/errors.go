package media

import apperrors "shared/app_errors"

var (
	ErrStorageNil      = apperrors.New(apperrors.Internal, "media service: storage is nil")
	ErrFileTooLarge    = apperrors.New(apperrors.MediaFileTooLarge, "file too large (max 5MB)")
	ErrInvalidFileType = apperrors.New(apperrors.MediaTypeInvalid, "only jpg, png, webp allowed")
	ErrInvalidURL      = apperrors.New(apperrors.MediaURLInvalid, "invalid URL")
	ErrEmptyURL        = apperrors.New(apperrors.MediaURLInvalid, "empty URL")
)
