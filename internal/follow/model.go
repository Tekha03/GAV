package follow

type Follow struct {
	FollowerID	uint	`gorm:"primaryKey"`
	FollowingID	uint	`gorm:"primaryKey"`
}

func NewFollow(followerID, followingID uint) *Follow {
	return &Follow{FollowerID: followerID, FollowingID: followingID}
}
