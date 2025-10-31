package user

import (
	"fmt"
	"time"
	"strings"
)

type UserSettings struct {
	private 		  bool
	emailNotification  bool
	password		  string
	createdAt		  time.Time
}

func NewUserSettings() *UserSettings {
	return &UserSettings{
		private: 			getPrivacy(),
		emailNotification:	 getEmailNotification(),
		password: 			getPassword(),
		createdAt: 			getCreationTime(),
	}
}

func getPrivacy() bool {
	var privacy string
	fmt.Println("Do you want your profile to be private?")
	fmt.Scanln(&privacy)

	for !correctInput(privacy) {
		fmt.Println("Please enter yes or no!")
		fmt.Scanln(&privacy)
	}

	if strings.ToLower(privacy) == "yes" {
		return true
	} else {
		return false
	}
}

func getEmailNotification() bool {
	var notifications string
	fmt.Println("Do you want to get notifications by email?")
	fmt.Scanln(&notifications)

	for !correctInput(notifications) {
		fmt.Println("Please enter yes or no!")
		fmt.Scanln(&notifications)
	}

	if strings.ToLower(notifications) == "yes" {
		return true
	} else {
		return false
	}
}

func correctInput(answer string) bool {
	if strings.ToLower(answer) == "yes" || strings.ToLower(answer) == "no" {
		return true
	} else {
		return false
	}
}

func getPassword() string {
	var password string
	fmt.Println("Please enter new password")
	fmt.Scanln(&password)

	for !isSafePassword(password) {
		fmt.Println("Your password does not meet security requirements, please enter safe password")
		safePasswordRequirements()
		fmt.Scanln(&password)
	}

	return password
}

func isSafePassword(password string) bool {
	// TODO
	return  true
}

func safePasswordRequirements() {
	// TODO
}

func getCreationTime() time.Time {
	return time.Time{}
}
