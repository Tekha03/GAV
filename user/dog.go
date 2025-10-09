package user

import (
	"fmt"
)

type Dog struct {
	dog_name	string
	age			uint
	breed		string
	character	string
	gender		string
}

func DogArrayConstructor() []*Dog {
	var quantity int
	fmt.Println("How many dogs do you have?")
	fmt.Scanln(&quantity)

	user_dogs := []*Dog{}

	for i := 0; i < quantity; i++ {
		dog := DogConstructor()
		user_dogs = append(user_dogs, dog)
	}

	return user_dogs
}

func DogConstructor() *Dog {
	return &Dog{
		dog_name: 	get_name(),
		age: 		get_age(),
		breed:		get_breed(),
		character: 	get_character(),
		gender: 	get_gender(),
	}
}

func get_name() string {
	fmt.Println("Please enter the name of your dog.")
	var dog_name string
	fmt.Scanln(&dog_name)

	return  dog_name
}

func get_age() uint {
	fmt.Println("How old is your dog?")
	var dog_age uint
	fmt.Scanln(&dog_age)

	return dog_age
}

func get_breed() string {
	fmt.Println("What breed is your dog?")
	var dog_breed string
	fmt.Scanln(&dog_breed)

	return dog_breed
}

func get_character() string {
	fmt.Println("What is your dog's personality like?")
	var dog_character string
	fmt.Scanln(&dog_character)

	return dog_character
}

func get_gender() string {
	fmt.Println("What is your dog's gender?")
	var dog_gender string
	fmt.Scanln(&dog_gender)

	return dog_gender
}
