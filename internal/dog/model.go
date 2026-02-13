package dog

import "github.com/google/uuid"

type Dog struct {
	ID 			uuid.UUID
	OwnerID		uuid.UUID

	Name		string
	Breed		string
	PhotoUrl	string

	Status		Status
	Age			Age
	Gender		Gender

	Lat      *float64
    Lon      *float64
    LocationVisible bool
}

func NewDog(
    ownerID uuid.UUID,
    name string,
    breed string,
    gender Gender,
    status Status,
    age Age,
    photoURL string,
) *Dog {
    return &Dog{
        OwnerID: ownerID,
        Name:    name,
        Breed:   breed,
        Gender:  gender,
        Status:  status,
        Age:     age,
        PhotoUrl: photoURL,
    }
}
