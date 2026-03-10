package dto

import "errors"

var (
	ErrPostNil 	= errors.New("post response: post is nil")
	ErrUserNil	= errors.New("user response: user is nil")
)
