package comment

import "time"

type Comment struct {
	ID			uint	`gorm:"primaryKey"`
	PostID		uint
	UserID		uint
	Content		string	`gorm:"not null"`
	CreatedAt	time.Time
}
