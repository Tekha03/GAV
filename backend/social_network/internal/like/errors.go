package like

import apperrors "shared/app_errors"

var (
	ErrRepoNil          = apperrors.New(apperrors.Internal, "like service: repo is nil")
	ErrInvalidLike      = apperrors.New(apperrors.Validation, "invalid like")
	ErrAlreadyLiked     = apperrors.New(apperrors.LikeAlreadyExists, "already liked")
	ErrLikeDoesNotExist = apperrors.New(apperrors.LikeNotFound, "like does not exist")
	ErrDBError          = apperrors.New(apperrors.Internal, "database error")
)
