package sqlite

import "errors"

var (
	ErrCommentNotFound      = errors.New("comment not found")
	ErrDogNotFound          = errors.New("dog not found")
	ErrSettingsNotFound     = errors.New("settings not found")
	ErrPostNotFound         = errors.New("post not found")
	ErrVaccinationNotFound  = errors.New("vaccination not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserExists           = errors.New("user already exists")
	ErrDBNil                = errors.New("base sqlite: db is nil")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")

	ErrMemberExists   = errors.New("chat member already exists")
	ErrMemberNotFound = errors.New("member not found")

	ErrChatExists   = errors.New("chat already exists")
	ErrChatNotFound = errors.New("chat not found")
	ErrNotGroup     = errors.New("operation allowed only for group chats")
	ErrEmptyTitle   = errors.New("title can not be empty")

	ErrMessageNotFound = errors.New("messsage not found")
	ErrMessageExists   = errors.New("message already exists")

	ErrAttachmentNotFound = errors.New("attachment not found")
	ErrAttachmentExist    = errors.New("attachment exist")

	ErrReactionExists   = errors.New("reaction already exists")
	ErrReactionNotFound = errors.New("reaction not found")
)
