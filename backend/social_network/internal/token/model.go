package token

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID			uuid.UUID	`gorm:"primaryKey"`
	UserID		uuid.UUID
	TokenHash	string
	ExpiresAt	time.Time
	Revoked		bool
	CreatedAt	time.Time
}
