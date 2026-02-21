package sqlite

import (
	"context"
	"gav/internal/follow"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FollowRepository struct {
	*BaseRepository
}

func NewFollowRepository(db *gorm.DB) follow.Repository {
	return &FollowRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *FollowRepository) Follow(ctx context.Context, follow follow.Follow) error {
	return r.DB(ctx).Create(&follow).Error
}

func (r *FollowRepository) Unfollow(ctx context.Context, unfollow follow.Follow) error {
	return r.DB(ctx).
	Where("follower_id = ? AND following_id = ?", unfollow.FollowerID, unfollow.FollowingID).
	Delete(&follow.Follow{}).Error
}

func (r *FollowRepository) FollowerExists(ctx context.Context, followCheck follow.Follow) (bool, error) {
	var count int64

	err := r.DB(ctx).Model(&follow.Follow{}).
	Where("follower_id = ? AND following_id = ?", followCheck.FollowerID, followCheck.FollowingID).
	Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, err
}

func (r *FollowRepository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]follow.Follow, error) {
	var followers []follow.Follow

	err := r.DB(ctx).Where("following_id = ?", userID).Find(&followers).Error
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (r *FollowRepository) GetFollowing(ctx context.Context, userID uuid.UUID) ([]follow.Follow, error) {
	var following []follow.Follow

	err := r.DB(ctx).Where("follower_id = ?", userID).Find(&following).Error
	if err != nil {
		return nil, err
	}

	return following, err
}
