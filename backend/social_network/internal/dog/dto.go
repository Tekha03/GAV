package dog

type CreateDogInput struct {
	Name		string
	Breed		string
	Gender		Gender
	Age			Age
	Status		Status
	PhotoUrl	string
}

type UpdateDogInput struct {
	Name		*string
	Breed		*string
	Gender		*Gender
	Age			*Age
	Status		*Status
	PhotoUrl	*string
}

type UpdateLocationInput struct {
    Latitude       float64
    Longitude      float64
    Status         LocationStatus
    ClearLocation  bool 
}