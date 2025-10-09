package user

type UserStats struct {
	posts			[]string
	followers		[]string
	following		[]string
	posts_cnt		uint
	followers_cnt	uint
	following_cnt	uint
}

func UserStatsConstructor() *UserStats {
	return  &UserStats{
		posts: 	   		[]string{},
		followers: 		[]string{},
		following: 		[]string{},
		posts_cnt: 	   	0,
		followers_cnt: 	0,
		following_cnt: 	0,
	}
}
