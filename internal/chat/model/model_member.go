package model

import (
	"time"

	"github.com/google/uuid"
)

type MemberRole string

const (
	Admin  MemberRole = "admin"
	Member MemberRole = "member"
)

type ChatMember struct {
	UserID   	  	  uuid.UUID
	ChatID    		  uuid.UUID
	Role      		  MemberRole
	JoinedAt 		  time.Time
	LeftAt   		  *time.Time
	Muted    		  bool
	LastReadMessageID uuid.UUID
}
