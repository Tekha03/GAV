package repository

import (
	"context"
	"gav/internal/chat/model"

	"github.com/google/uuid"
)

type ChatMemberRepository interface {
	AddMember(ctx context.Context, member *model.ChatMember) error
	RemoveMember(ctx context.Context, memberID, chatID uuid.UUID) error
	GetMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error)
	UpdateRole(ctx context.Context, chatID, userID uuid.UUID, role *model.MemberRole) error
	SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}
