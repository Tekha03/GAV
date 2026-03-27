package service

import (
	"context"
	"messanger/internal/model"

	"github.com/google/uuid"
)

type Service interface {
	CreatePrivateChat(ctx context.Context, userID1, userID2 uuid.UUID) (*model.Chat, error)
	CreateGroupChat(ctx context.Context, title string, creatorID uuid.UUID, membersIDs []uuid.UUID) (*model.Chat, error)
	GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error)

	AddMember(ctx context.Context, userID, chatID uuid.UUID) error
	RemoveMember(ctx context.Context, userID, chatID uuid.UUID) error
	GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error)
	LeaveChat(ctx context.Context, userID, chatID uuid.UUID) error
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]*model.Chat, error)

	UpdateChatTitle(ctx context.Context, chatID uuid.UUID, newTitle string) error
	UpdateChatPhoto(ctx context.Context, chatID uuid.UUID, newPhotoURL string) error

	SendMessage(ctx context.Context, input model.SendMessageInput) (*model.Message, error)
	EditMessage(ctx context.Context, messageID uuid.UUID, newText string) (*model.Message, error)
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
	GetMessages(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error)
	MarkAsRead(ctx context.Context, chatID, userID uuid.UUID) error

	PinMessage(ctx context.Context, messageID uuid.UUID) error
	UnpinMessage(ctx context.Context, messageID uuid.UUID) error
	GetPinnedMessages(ctx context.Context, chatID uuid.UUID) ([]*model.Message, error)

	ForwardMessage(ctx context.Context, messageID, targetChatID, senderID uuid.UUID) (*model.Message, error)

	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetChatUnreadCount(ctx context.Context, chatID, userID uuid.UUID) (int, error)

	SendTyping(ctx context.Context, chatID, userID uuid.UUID) error

	AddReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error
	RemoveReaction(ctx context.Context, messageID, userID uuid.UUID) error
}