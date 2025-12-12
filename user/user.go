package user

type User struct {
	ID			uint
	Profile		 UserProfile
	Stats		UserStats
	Settings	UserSettings
	UserDogs 	[]*Dog
}

func NewUser(id uint, profile *UserProfile, stats *UserStats, settings *UserSettings, dogs []*Dog) *User {
	var prof UserProfile
	var st UserStats
	var set UserSettings

	if profile != nil {
		prof = *profile
	}

	if stats != nil {
		st = *stats
	}

	if settings != nil {
		set = *settings
	}

	return &User{
		ID: 		id,
		Profile: 	 prof,
		Stats:		st,
		Settings: 	set,
		UserDogs:	dogs,
	}
}
