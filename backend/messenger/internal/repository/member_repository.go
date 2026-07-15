package repository

import (
	"context"
	"messenger/internal/model"

	"github.com/google/uuid"
)

type ChatMemberRepository interface {
	AddMember(ctx context.Context, member *model.ChatMember) error
	RemoveMember(ctx context.Context, memberID, chatID uuid.UUID) error
	GetMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error)
	UpdateRole(ctx context.Context, chatID, userID uuid.UUID, role *model.MemberRole) error
	SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	FindPrivateChatBetween(ctx context.Context, userID1, userID2 uuid.UUID) (uuid.UUID, error)
	GetLastReadMessageID(ctx context.Context, chatID, userID uuid.UUID) (uuid.UUID, error)
	MemberExists(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (bool, error)
}
