package user

import (
	"time"
)

type User struct {
	id			uint
	profile		 UserProfile
	stats		UserStats
	settings	UserSettings
	user_dogs 	[]Dog
}

type UserProfile struct {
	name	     string
	surname      string
	username	 string
	email		 string
	birth_date 	 string
	profile_photo string
	address		 string
	bio			 string
}

type Dog struct {
	dog_name	string
	age			uint
	breed		string
	character	string
	gender		string
}

type UserStats struct {
	posts			[]string
	followers		[]string
	following		[]string
	posts_cnt		uint
	followers_cnt	uint
	following_cnt	uint
}

type UserSettings struct {
	private 		  bool
	email_notification bool
	password		  string
	createdAt		  time.Time
}

func (user *User) UserConstructor(id uint) *User {
	return &User{
		id: 		id,
		profile: 	 *user.profile.UserProfileConstructor("", "", ""),
		stats:		*user.stats.UserStatsConstructor(),
		settings: 	*user.settings.UserSettingsConstructor("", time.Time{}),
		user_dogs: []Dog{},
	}
}

func (user_profile *UserProfile) UserProfileConstructor(name string, username string, birth_date string) *UserProfile {
	return &UserProfile{
		name: 			name,
		surname: 		"",
		username: 		username,
		email: 			"",
		birth_date: 	birth_date,
		profile_photo: 	 "",
		address: 		"",
		bio: 			"",
	}
}

func (dog *Dog) DogConstructor(dog_name string, gender string) *Dog {
	return &Dog{
		dog_name: 	dog_name,
		age: 		0,
		breed: 		"",
		character: 	"",
		gender: 	gender,
	}
}

func (user_stats *UserStats) UserStatsConstructor() *UserStats {
	return  &UserStats{
		posts: 	   		[]string{},
		followers: 		[]string{},
		following: 		[]string{},
		posts_cnt: 	   	0,
		followers_cnt: 	0,
		following_cnt: 	0,
	}
}

func (user_settings *UserSettings) UserSettingsConstructor(password string, createdAt time.Time) *UserSettings {
	return &UserSettings{
		private: 			false,
		email_notification:	 true,
		password: 			password,
		createdAt: 			createdAt,
	}
}
