package migrations

import (
	"gav/internal/comment"
	"gav/internal/follow"
	"gav/internal/like"
	"gav/internal/role"
	"gav/internal/user"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.User{},
		&role.Role{},
		&follow.Follow{},
		&like.Like{},
		&comment.Comment{},
	)
}
