package user

type User struct {
	id			  int
	name 		  string
	surname 	  string
	email		  string
	password	  string
	birth_date 	  string
	dog_breed 	  string
	dog_name	  string
	profile_photo  string
	address		  string
}

func (user *User) GetID() int {
	return  user.id
}

func (user *User) GetPassword() string {
	return  user.password
}
