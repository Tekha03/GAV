package sqlite

import "errors"

var (
	ErrCommentNotFound		= errors.New("comment not found")
	ErrDogNotFound 			= errors.New("dog not found")
	ErrSettingsNotFound 	= errors.New("settings not found")
	ErrPostNotFound 		= errors.New("post not found")
	ErrVaccinationNotFound 	= errors.New("vaccination not found")
	ErrUserNotFound 		= errors.New("user not found")
	ErrUserExists 			= errors.New("user already exists")
	ErrDBNil				= errors.New("base sqlite: db is nil")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
