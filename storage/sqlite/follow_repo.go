package sqlite

import (
	"gav/internal/follow"

	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (fr *FollowRepository) Follow(followerID, followingID uint) error {
	return fr.db.Create(&follow.Follow{
		FollowerID: followerID,
		FollowingID: followingID,
	}).Error
}
