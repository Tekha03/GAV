package follow

import "errors"

var (
	ErrCannotFollowYourself = errors.New("you cannot follow yourself.")
	ErrAlreadyFollowing 	= errors.New("already following")
	ErrInvalidUserID 		= errors.New("invalid user id")

	ErrFollowerIDNil 		= errors.New("follow model: follower id is nil")
	ErrFollowingIDNil 		= errors.New("follow model: following id is nil")
	ErrRepoNil				= errors.New("follow service: repo is nil")
)
