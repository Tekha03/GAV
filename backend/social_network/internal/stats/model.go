package stats

import (
	"time"

	"github.com/google/uuid"
)

type UserStats struct {
	UserID		uuid.UUID	`gorm:"primaryKey"`
	PostCount	uint
	Followers	uint
	Followings	uint
	DogsCount	uint
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

type PostStats struct {
	PostID			uuid.UUID	`gorm:"primaryKey"`
	AuthorID		uuid.UUID	`gorm:"index"`
	LikesCount		uint
	CommentsCount	uint
	CreatedAt		time.Time
	UpdatedAt		time.Time
}

type ProfileStats struct {
	UserID		uuid.UUID
	PostCount	uint
	Followers	uint
	Followings	uint
}
