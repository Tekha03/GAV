package user

type UserStats struct {
	posts			[]string
	followers		[]string
	following		[]string
	postsCnt		uint
	followersCnt	uint
	followingCnt	uint
}

func NewUserStats() *UserStats {
	return  &UserStats{
		posts: 	   		[]string{},
		followers: 		[]string{},
		following: 		[]string{},
		postsCnt: 	   	0,
		followersCnt: 	0,
		followingCnt: 	0,
	}
}
