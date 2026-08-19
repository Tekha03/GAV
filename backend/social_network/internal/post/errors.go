package post

import apperrors "shared/app_errors"

var (
	ErrPostNotFound = apperrors.New(apperrors.PostNotFound, "post not found")
	ErrForbidden    = apperrors.New(apperrors.PostAccessDenied, "forbidden")
	ErrEmptyContent = apperrors.New(apperrors.PostContentRequired, "post content is required")

	ErrRepoNil = apperrors.New(apperrors.Internal, "post service: repo is nil")
)
