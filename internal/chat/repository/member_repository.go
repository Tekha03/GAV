package repository

import (
	"context"
	"gav/internal/chat"

	"github.com/google/uuid"
)

type ChatMemberRepository interface {
	AddMember(ctx context.Context, member *chat.ChatMember) error
	RemoveMember(ctx context.Context, memberID, chatID uuid.UUID) error
	GetMembers(ctx context.Context, chatID uuid.UUID) ([]*chat.ChatMember, error)
	UpdateRole(ctx context.Context, chatID, userID uuid.UUID, role *chat.MemberRole) error
	SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error
}
