package user

type UpdateUserInput struct {
	Email		*string
	Password	*string
	Role		*string
}

type UpdateLocationInput struct {
    Latitude       float64
    Longitude      float64
    Status         LocationStatus
    ClearLocation  bool
}
