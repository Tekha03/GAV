package dto

import (
	"messenger/internal/model"
	"time"

	"github.com/google/uuid"
)

type ChatMemberRequestDTO struct {
	UserID uuid.UUID        `json:"user_id"`
	ChatID uuid.UUID        `json:"chat_id"`
	Role   model.MemberRole `json:"role"`
	Muted  bool             `json:"muted"`
}

type ChatMemberResponseDTO struct {
	UserID            uuid.UUID        `json:"user_id"`
	ChatID            uuid.UUID        `json:"chat_id"`
	Role              model.MemberRole `json:"role"`
	JoinedAt          time.Time        `json:"joined_at"`
	LeftAt            *time.Time       `json:"left_at,omitempty"`
	Muted             bool             `json:"muted"`
	LastReadMessageID uuid.UUID        `json:"last_read_message_id"`
}
