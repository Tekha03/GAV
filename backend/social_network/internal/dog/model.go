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

	Lat         *float64
    Lon         *float64
    LocationStatus  LocationStatus
    Visibility      LocationVisibility 
}

func NewDog(
    ownerID uuid.UUID,
    name string,
    breed string,
    gender Gender,
    status Status,
    age Age,
    photoURL string,
) (*Dog, error) {

    if ownerID == uuid.Nil {
        return nil, ErrOwnerIDNil
    }
    if name == "" {
        return nil, ErrNameEmpty
    }
    if breed == "" {
        return nil, ErrBreedEmpty
    }
    if gender == "" {
        return nil, ErrGenderEmpty
    }
    if status == "" {
        return nil, ErrStatusEmpty
    }
    if age == "" {
        return nil, ErrAgeEmpty
    }
    if photoURL == "" {
        return nil, ErrPhotoURLEmpty
    }

    return &Dog{
        OwnerID: ownerID,
        Name:    name,
        Breed:   breed,
        Gender:  gender,
        Status:  status,
        Age:     age,
        PhotoUrl: photoURL,
    }, nil
}
