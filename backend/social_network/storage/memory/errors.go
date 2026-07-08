package memory

import "errors"

var (
	ErrCommentNotFound     = errors.New("comment not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExists          = errors.New("user already exists")
	ErrVaccinationExists   = errors.New("vaccination already exists")
	ErrVaccinationNotFound = errors.New("vaccination not found")
	ErrStatExist           = errors.New("stat exist in repository")
	ErrStatNotFound        = errors.New("stat not found")

	ErrCommentNil     = errors.New("comment memory: comment is nil")
	ErrDogNil         = errors.New("dog memory: dog is nil")
	ErrPostNil        = errors.New("post memory: post is nil")
	ErrStatsNil       = errors.New("stats memory: stats is nil")
	ErrUserNil        = errors.New("user memory: user is nil")
	ErrVaccinationNil = errors.New("vaccination memory: vaccination is nil")

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

	ErrDogNotFound = errors.New("dog not found")
	ErrDogExists   = errors.New("dog exists in repository")
)
