package dog

type Dog struct {
	ID 			uint
	OwnerID		uint

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
    ownerID uint,
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
