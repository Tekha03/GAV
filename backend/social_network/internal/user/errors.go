package user

import "errors"

var (
	ErrEmailEmpty        = errors.New("user model: email is empty")
	ErrPasswordHashEmpty = errors.New("user mode: password hash is empty")
	ErrRepoNil           = errors.New("user service: repo is nil")
	ErrUserNotFound      = errors.New("user not found")
	ErrFail              = errors.New("fail")
)
