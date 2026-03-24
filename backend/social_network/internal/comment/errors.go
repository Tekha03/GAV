package comment

import "errors"

var (
	ErrRepoEmpty 	= errors.New("comment service: repo is nil")
	ErrDB			= errors.New("db error")
)