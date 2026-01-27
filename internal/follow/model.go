package follow

type Follow struct {
	FollowerID	uint	`gorm:"primaryKey"`
	FollowingID	uint	`gorm:"primaryKey"`
}
