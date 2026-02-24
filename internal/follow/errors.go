package follow

import "errors"

var (
	ErrFollowerIDNil 	= errors.New("follow model: follower id is nil")
	ErrFollowingIDNil 	= errors.New("follow model: following id is nil")
	ErrRepoNil			= errors.New("follow service: repo is nil")
)
