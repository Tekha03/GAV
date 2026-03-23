package handlers

import "errors"


var (
	ErrInvalidInput 	= errors.New("invalid input")
	ErrUnauthorized 	= errors.New("unauthorized")

	ErrAuthNil			= errors.New("auth handler: service is nil")
	ErrCommentNil		= errors.New("comment handler: service is nil")
	ErrDogNil			= errors.New("dog handler: service is nil")
	ErrFeedNil			= errors.New("feed handler: service is nil")
	ErrFollowNil		= errors.New("follow handler: service is nil")
	ErrLikeNil			= errors.New("like handler: service is nil")
	ErrPostNil			= errors.New("post handler: service is nil")
	ErrProfileNil		 = errors.New("profile handler: service is nil")
	ErrSettingsNil		= errors.New("settins handler: service is nil")
	ErrStatsNil		 	= errors.New("stats handler: service is nil")
	ErrUserNil		 	= errors.New("user handler: service is nil")
	ErrVaccinationNil	= errors.New("vaccination handler: service is nil")
	ErrMediaNil			= errors.New("upload handler: service is nil")
)
