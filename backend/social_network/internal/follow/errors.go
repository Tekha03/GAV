package follow

import apperrors "shared/app_errors"

var (
	ErrCannotFollowYourself = apperrors.New(apperrors.FollowSelfNotAllowed, "you cannot follow yourself")
	ErrAlreadyFollowing     = apperrors.New(apperrors.FollowAlreadyExists, "already following")
	ErrInvalidUserID        = apperrors.New(apperrors.Validation, "invalid user id", apperrors.WithDetail("field", "user_id"))

	ErrFollowerIDNil  = apperrors.New(apperrors.Validation, "follow model: follower id is nil")
	ErrFollowingIDNil = apperrors.New(apperrors.Validation, "follow model: following id is nil")
	ErrRepoNil        = apperrors.New(apperrors.Internal, "follow service: repo is nil")
	ErrDBError        = apperrors.New(apperrors.Internal, "database error")
)
