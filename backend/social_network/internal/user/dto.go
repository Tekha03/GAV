package user

type UpdateUserInput struct {
	Email    *string
	Password *string
	Role     *string
}

type UpdateLocationInput struct {
	Latitude      float64            `json:"lat"`
	Longitude     float64            `json:"lon"`
	Status        LocationStatus     `json:"location_status"`
	Visibility    LocationVisibility `json:"visibility"`
	ClearLocation bool               `json:"clear_location"`
}
