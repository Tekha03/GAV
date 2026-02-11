package chat

import "time"

type MemberRole string

const (
	Admin MemberRole = "admin"
	Member MemberRole = "member"
)

type ChatMember struct {
	UserID 		uint
	ChatID		uint
	Role		MemberRole
	JoinedAt	time.Time
	LeftAt		*time.Time
	Muted		bool
}