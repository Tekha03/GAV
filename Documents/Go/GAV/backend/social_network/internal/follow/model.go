package follow

import "github.com/google/uuid"

type Follow struct {
	FollowerID	uuid.UUID	`gorm:"primaryKey"`
	FollowingID	uuid.UUID	`gorm:"primaryKey"`
}

func NewFollow(followerID, followingID uuid.UUID) (*Follow, error) {
	if followerID == uuid.Nil {
		return nil, ErrFollowerIDNil
	}
	if followingID == uuid.Nil {
		return nil, ErrFollowingIDNil
	}

	return &Follow{FollowerID: followerID, FollowingID: followingID}, nil
}
