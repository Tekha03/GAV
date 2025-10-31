package user

type User struct {
	id			uint
	profile		 UserProfile
	stats		UserStats
	settings	UserSettings
	userDogs 	[]*Dog
}

func (user *User) UserConstructor(id uint) *User {
	return &User{
		id: 		id,
		profile: 	 *NewUserProfile("", "", ""),
		stats:		*NewUserStats(),
		settings: 	*NewUserSettings(),
		userDogs:	NewDogArray(),
	}
}
