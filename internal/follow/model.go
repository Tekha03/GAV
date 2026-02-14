package follow

import "github.com/google/uuid"

type Follow struct {
	FollowerID	uuid.UUID	`gorm:"primaryKey"`
	FollowingID	uuid.UUID	`gorm:"primaryKey"`
}

func NewFollow(followerID, followingID uuid.UUID) *Follow {
	return &Follow{FollowerID: followerID, FollowingID: followingID}
}
