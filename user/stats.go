package user

type UserStats struct {
	Posts			[]string	`gorm:"type:json"`
	Followers		[]string	`gorm:"type:json"`
	Following		[]string	`gorm:"type:json"`
	PostsCnt		uint
	FollowersCnt	uint
	FollowingCnt	uint
}

func NewUserStats() *UserStats {
	return  &UserStats{
		Posts: 	   		[]string{},
		Followers: 		[]string{},
		Following: 		[]string{},
		PostsCnt: 	   	0,
		FollowersCnt: 	0,
		FollowingCnt: 	0,
	}
}
