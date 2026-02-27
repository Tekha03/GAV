package token

import "errors"

var (
	ErrRepoNil					= errors.New("token service: repo is nil")
	ErrInvalidRefresh 			= errors.New("token service: invalid refresh")
)
