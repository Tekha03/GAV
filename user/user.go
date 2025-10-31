package user

type User struct {
	id			uint
	profile		 UserProfile
	stats		UserStats
	settings	UserSettings
	userDogs 	[]*Dog
}

func (user *User) UserConstructor(id uint) *User {
	profile := NewUserProfile("", "", "")
	stats := NewUserStats()
	settings := NewUserSettings()
	userDogs := NewDogArray()

	return &User{
		id: 		id,
		profile: 	 *profile,
		stats:		*stats,
		settings: 	*settings,
		userDogs:	userDogs,
	}
}
