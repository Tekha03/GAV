package user

type User struct {
	id			uint
	profile		 UserProfile
	stats		UserStats
	settings	UserSettings
	user_dogs 	[]*Dog
}

func (user *User) UserConstructor(id uint) *User {
	profile := UserProfileConstructor("", "", "")
	stats := UserStatsConstructor()
	settings := UserSettingsConstructor()
	user_dogs := DogArrayConstructor()

	return &User{
		id: 		id,
		profile: 	 *profile,
		stats:		*stats,
		settings: 	*settings,
		user_dogs:	user_dogs,
	}
}
