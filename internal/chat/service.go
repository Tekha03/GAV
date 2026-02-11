package chat

import "context"

type Service interface {
	CreatePrivateChat(ctx context.Context, userID1, userID2 uint) error
	CreateGroupChat(ctx context.Context, title string, creatorID uint, membersIDs []uint) error

	AddMember(ctx context.Context, userID uint) error
	RemoveMember(ctx context.Context, userID uint) error
	LeaveChat(ctx context.Context, userID, chatID uint) error
	GetUserChats(ctx context.Context, userID uint) ([]*Chat, error)

	SendMessage(ctx context.Context, input SendMessageInput) (*Message, error)
	EditMessage(ctx context.Context, newText string) (*Message, error)
	DeleteMessage(ctx context.Context, messageID uint) error
	GetMessages(ctx context.Context, chatID, limit uint, cursorID *uint) ([]*Message, error)

	AddReaction(ctx context.Context, messageID, userID uint, emoji string) error
	RemoveReaction(ctx context.Context, messageID, userID uint) error
}