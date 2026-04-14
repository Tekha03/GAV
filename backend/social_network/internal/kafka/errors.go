package kafka

import (
	"errors"
)

var (
	ErrMessageUseCaseNil = errors.New("message use case is nil")
	ErrChatUseCaseNil = errors.New("chat use case is nil")
	ErrReactionUseCaseNil = errors.New("reaction use case is nil")
)
