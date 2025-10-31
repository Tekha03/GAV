package user

type UserProfile struct {
	name	     string
	surname      string
	username	 string
	email		 string
	birthDate 	 string
	profilePhoto string
	address		 string
	bio			 string
}

func NewUserProfile(name string, username string, birthDate string) *UserProfile {
	return &UserProfile{
		name: 			name,
		surname: 		"",
		username: 		username,
		email: 			"",
		birthDate: 		birthDate,
		profilePhoto: 	 "",
		address: 		"",
		bio: 			"",
	}
}
