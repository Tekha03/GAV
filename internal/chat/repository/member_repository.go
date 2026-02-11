package repository

import (
	"context"
	"gav/internal/chat"
)

type ChatMemberRepository interface {
	AddMember(ctx context.Context, member *chat.ChatMember) error
	RemoveMember(ctx context.Context, memberID uint) error
	GetMembers(ctx context.Context, chatID uint) ([]*chat.ChatMember, error)
	GetUserChats(ctx context.Context, userID uint) ([]*chat.ChatMember, error)
	UpdateRole(ctx context.Context, chatID, userID uint, role chat.MemberRole) error
	SetMuted(ctx context.Context, chatID, userID uint, muted bool) error
}
