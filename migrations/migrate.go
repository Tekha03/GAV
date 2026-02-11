package migrations

import (
	"gav/internal/comment"
	"gav/internal/follow"
	"gav/internal/like"
	"gav/internal/user"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.User{},
		&follow.Follow{},
		&like.Like{},
		&comment.Comment{},
	)
}
