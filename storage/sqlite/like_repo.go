package sqlite

import (
	"gav/internal/like"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (lr *LikeRepository) Add(userID, postID uint) error {
	return lr.db.Create(&like.Like{UserID: userID, PostID: postID}).Error
}
