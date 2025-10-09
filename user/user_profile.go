package user

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

func UserProfileConstructor(name string, username string, birth_date string) *UserProfile {
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
