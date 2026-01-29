package dog

type Dog struct {
	ID 			int
	OwnerID		int

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