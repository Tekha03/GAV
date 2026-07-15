package errors

import "errors"

var (
	ErrChatWithYourSelf      = errors.New("can not create private chat with yourself of add yourself to group")
	ErrChatNotFound          = errors.New("chat not found")
	ErrNoMembers             = errors.New("no members in chat")
	ErrNoChats               = errors.New("no chats")
	ErrTitleUpdate           = errors.New("error to update chat title")
	ErrPhotoUpdate           = errors.New("error to update chat photo")
	ErrEmptyMessage          = errors.New("message content can not be empty")
	ErrMessageNotFound       = errors.New("message not found")
	ErrInvalidReply          = errors.New("error to reply. Message from other chat")
	ErrTextOverLength        = errors.New("excessive message length")
	ErrAttachmentsOverLength = errors.New("excessive attacments count")
	ErrIsNotMember           = errors.New("this user is not member of chat")
	ErrMemberExists          = errors.New("chat member already exists")
	ErrMemberNotFound        = errors.New("member not found")
	ErrChatAccessDenied      = errors.New("chat access denied")
)
