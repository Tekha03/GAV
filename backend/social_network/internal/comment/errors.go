package comment

import apperrors "shared/app_errors"

var (
	ErrRepoEmpty = apperrors.New(apperrors.Internal, "comment service: repo is nil")
	ErrDB        = apperrors.New(apperrors.Internal, "database error")
)
