package token

import apperrors "shared/app_errors"

var (
	ErrRepoNil        = apperrors.New(apperrors.Internal, "token service: repo is nil")
	ErrInvalidRefresh = apperrors.New(apperrors.AuthRefreshInvalid, "invalid refresh token")
	ErrTokenNotFound  = apperrors.New(apperrors.AuthRefreshInvalid, "token not found")
	ErrFail           = apperrors.New(apperrors.Internal, "fail")
)
