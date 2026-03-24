package comment

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID			uuid.UUID	`gorm:"primaryKey"`
	PostID		uuid.UUID
	UserID		uuid.UUID
	Content		string		`gorm:"not null"`
	CreatedAt	time.Time
}
