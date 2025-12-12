package user

type UserProfile struct {
	Name	     string
	Surname      string
	Username	 string
	Email		 string
	BirthDate 	 string
	ProfilePhoto  string
	Address		 string
	Bio			 string
}

func NewUserProfile(name, surname, username, email, birthDate, profilePhoto, address, bio string) *UserProfile {
	return &UserProfile{
		Name: 			name,
		Surname: 		surname,
		Username: 		username,
		Email: 			email,
		BirthDate: 		birthDate,
		ProfilePhoto: 	 profilePhoto,
		Address: 		address,
		Bio: 			bio,
	}
}

// func getName() string {
// 	fmt.Println("What's your name?")
// 	var name string
// 	fmt.Scanln(&name)

// 	return name
// }

// func getSurname() string {
// 	fmt.Println("What's your surname?")
// 	var surname string
// 	fmt.Scanln(&surname)

// 	return surname
// }

// func getUsername() string {
// 	fmt.Println("Enter your username")
// 	var username string
// 	fmt.Scanln(&username)

// 	return username
// }

// func getEmail() string {
// 	fmt.Println("Enter your email")
// 	var email string
// 	fmt.Scanln(&email)

// 	return email
// }

// func getBirthDate() string {
// 	fmt.Println("Enter your birth date")
// 	var birthDate string
// 	fmt.Scanln(&birthDate)

// 	return birthDate
// }

// func getProfilePhoto() string {
// 	fmt.Println("Upload your profile photo")
// 	var profilePhoto string
// 	fmt.Scanln(&profilePhoto)

// 	return profilePhoto
// }

// func getAddress() string {
// 	fmt.Println("Enter your address")
// 	var address string
// 	fmt.Scanln(&address)

// 	return address
// }

// func getBio() string {
// 	fmt.Println("Tell something about yourself")
// 	var bio string
// 	fmt.Scanln(&bio)

// 	return bio
// }
