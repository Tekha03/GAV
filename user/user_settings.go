package user

import (
	"fmt"
	"time"
	"strings"
)

type UserSettings struct {
	private 		  bool
	email_notification bool
	password		  string
	createdAt		  time.Time
}

func UserSettingsConstructor() *UserSettings {
	return &UserSettings{
		private: 			get_privacy(),
		email_notification:	 get_email_norification(),
		password: 			get_password(),
		createdAt: 			get_creation_time(),
	}
}

func get_privacy() bool {
	var privacy string
	fmt.Println("Do you want your profile to be private?")
	fmt.Scanln(&privacy)

	for !correct_input(privacy) {
		fmt.Println("Please enter yes or no!")
		fmt.Scanln(&privacy)
	}

	if strings.ToLower(privacy) == "yes" {
		return true
	} else {
		return false
	}
}

func get_email_norification() bool {
	var notifications string
	fmt.Println("Do you want to get notifications by email?")
	fmt.Scanln(&notifications)

	for !correct_input(notifications) {
		fmt.Println("Please enter yes or no!")
		fmt.Scanln(&notifications)
	}

	if strings.ToLower(notifications) == "yes" {
		return true
	} else {
		return false
	}
}

func correct_input(answer string) bool {
	if strings.ToLower(answer) == "yes" || strings.ToLower(answer) == "no" {
		return true
	} else {
		return false
	}
}

func get_password() string {
	var password string
	fmt.Println("Please enter new password")
	fmt.Scanln(&password)

	for !is_safe_password(password) {
		fmt.Println("Your password does not meet security requirements, please enter safe password")
		safe_password_requirements()
		fmt.Scanln(&password)
	}

	return password
}

func is_safe_password(password string) bool {
	// TODO
	return  true
}

func safe_password_requirements() {
	// TODO
}

func get_creation_time() time.Time {
	return time.Time{}
}
