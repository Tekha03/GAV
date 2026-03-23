package like

import "errors"

var (
	ErrRepoNil 			= errors.New("like service: repo is nil")
	ErrInvalidLike 		= errors.New("invalid like")
	ErrAlreadyLiked 	= errors.New("already liked")
	ErrLikeDoesNotExist = errors.New("like does not exist")
	ErrDBError			= errors.New("db error")
)