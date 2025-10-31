package user

import (
	"fmt"
)

type Dog struct {
	dogName		string
	age			uint
	breed		string
	character	string
	gender		string
}

func NewDogArray() []*Dog {
	var quantity int
	fmt.Println("How many dogs do you have?")
	fmt.Scanln(&quantity)

	userDogs := []*Dog{}

	for i := 0; i < quantity; i++ {
		dog := NewDog()
		userDogs = append(userDogs, dog)
	}

	return userDogs
}

func NewDog() *Dog {
	return &Dog{
		dogName: 	getName(),
		age: 		getAge(),
		breed:		getBreed(),
		character: 	getCharacter(),
		gender: 	getGender(),
	}
}

func getName() string {
	fmt.Println("Please enter the name of your dog.")
	var dogName string
	fmt.Scanln(&dogName)

	return  dogName
}

func getAge() uint {
	fmt.Println("How old is your dog?")
	var dogAge uint
	fmt.Scanln(&dogAge)

	return dogAge
}

func getBreed() string {
	fmt.Println("What breed is your dog?")
	var dogBreed string
	fmt.Scanln(&dogBreed)

	return dogBreed
}

func getCharacter() string {
	fmt.Println("What is your dog's personality like?")
	var dogCharacter string
	fmt.Scanln(&dogCharacter)

	return dogCharacter
}

func getGender() string {
	fmt.Println("What is your dog's gender?")
	var dogGender string
	fmt.Scanln(&dogGender)

	return dogGender
}
