package memory

import "errors"

var (
	ErrCommentNotFound 		= errors.New("comment not found")
	ErrUserNotFound 		= errors.New("user not found")
	ErrUserExists 			= errors.New("user already exists")
	ErrVaccinationExists 	= errors.New("vaccination already exists")
	ErrVaccinationNotFound 	= errors.New("vaccination not found")
	ErrStatExist 			= errors.New("stat exist in repository")
	ErrStatNotFound			= errors.New("stat not found")

	ErrCommentNil 			= errors.New("comment memory: comment is nil")
	ErrDogNil				= errors.New("dog memory: dog is nil")
	ErrPostNil				= errors.New("post memory: post is nil")
	ErrStatsNil				= errors.New("stats memory: stats is nil")
	ErrUserNil				= errors.New("user memory: user is nil")
	ErrVaccinationNil		= errors.New("vaccination memory: vaccination is nil")

	ErrMemberExists = errors.New("chat member already exists")
	ErrMemberNotFound = errors.New("member not found")
)
