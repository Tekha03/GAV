package stats

import (
	"time"

	"github.com/google/uuid"
)

type UserStats struct {
	UserID     uuid.UUID `gorm:"primaryKey" json:"user_id"`
	PostCount  uint      `json:"post_count"`
	Followers  uint      `json:"followers_count"`
	Followings uint      `json:"followings_count"`
	DogsCount  uint      `json:"dogs_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostStats struct {
	PostID        uuid.UUID `gorm:"primaryKey" json:"post_id"`
	AuthorID      uuid.UUID `gorm:"index" json:"author_id"`
	LikesCount    uint      `json:"likes_count"`
	CommentsCount uint      `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProfileStats struct {
	UserID     uuid.UUID `json:"user_id"`
	PostCount  uint      `json:"post_count"`
	Followers  uint      `json:"followers_count"`
	Followings uint      `json:"followings_count"`
}
