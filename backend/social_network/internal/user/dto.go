package user

type UpdateUserInput struct {
	Email *string `json:"email"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangeRoleInput struct {
	Role string `json:"role"`
}

type UpdateLocationInput struct {
	Latitude      float64            `json:"lat"`
	Longitude     float64            `json:"lon"`
	Status        LocationStatus     `json:"location_status"`
	Visibility    LocationVisibility `json:"visibility"`
	ClearLocation bool               `json:"clear_location"`
}
