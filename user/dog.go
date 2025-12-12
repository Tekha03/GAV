package user

type Dog struct {
	DogName		string
	Age			uint
	Breed		string
	Character	string
	Gender		string
}

func NewDog(dogName string, age uint, breed, character, gender string) *Dog {
	return &Dog{
		DogName: 	dogName,
		Age: 		age,
		Breed:		breed,
		Character: 	character,
		Gender: 	gender,
	}
}

// func DogName() string {
// 	fmt.Println("Please enter the name of your dog.")
// 	var dogName string
// 	fmt.Fscanln(Input, &dogName)

// 	return  dogName
// }

// func Age() uint {
// 	fmt.Println("How old is your dog?")
// 	var dogAge uint
// 	fmt.Fscanln(Input, &dogAge)

// 	return dogAge
// }

// func Breed() string {
// 	fmt.Println("What breed is your dog?")
// 	var dogBreed string
// 	fmt.Fscanln(Input, &dogBreed)

// 	return dogBreed
// }

// func Character() string {
// 	fmt.Println("What is your dog's personality like?")
// 	var dogCharacter string
// 	fmt.Fscanln(Input, &dogCharacter)

// 	return dogCharacter
// }

// func Gender() string {
// 	fmt.Println("What is your dog's gender?")
// 	var dogGender string
// 	fmt.Fscanln(Input, &dogGender)

// 	return dogGender
// }
