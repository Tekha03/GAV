package repository

import "errors"

var (
	ErrAttachmentNotFound = errors.New("attachment not found")
	ErrChatExists         = errors.New("chat already exists")
	ErrChatNotFound       = errors.New("chat not found")
	ErrMemberExists       = errors.New("member already exists")
	ErrMemberNotFound     = errors.New("member not found")
	ErrMessageExists      = errors.New("message already exists")
	ErrMessageNotFound    = errors.New("message not found")
	ErrReactionExists     = errors.New("reaction already exists")
	ErrReactionNotFound   = errors.New("reaction not found")
)
