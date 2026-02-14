package sqlite

import (
	"context"

	"gav/internal/like"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (lr *LikeRepository) Add(ctx context.Context, like like.Like) error {
	exists, err := lr.LikeExists(ctx, like)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return lr.db.WithContext(ctx).Create(&like).Error
}

func (lr *LikeRepository) Remove(ctx context.Context, likeToRemove like.Like) error {
	result := lr.db.WithContext(ctx).Where(
		"user_id = ? AND post_id = ?", likeToRemove.UserID, likeToRemove.PostID,
	).Delete(&like.Like{})

	return result.Error
}

func (lr *LikeRepository) LikeExists(ctx context.Context, likeToCheck like.Like) (bool, error) {
	var count int64
	err := lr.db.WithContext(ctx).Model(&like.Like{}).Where(
		"user_id = ? AND post_id = ?", likeToCheck.UserID, likeToCheck.PostID,
	).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
