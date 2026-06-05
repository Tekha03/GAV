package profile

type CreateProfileInput struct {
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Username        string `json:"username" validate:"required,min=3,max=30"`
	ProfilePhotoUrl string `json:"profile_photo_url"`
	Bio             string `json:"bio"`
	Address         string `json:"address"`
	BirthDate       string `json:"birth_date"`
}

type UpdateProfileInput struct {
	Name            *string `json:"name"`
	Surname         *string `json:"surname"`
	Username        *string `json:"username" validate:"omitempty,min=3,max=30"`
	ProfilePhotoUrl *string `json:"profile_photo_url"`
	Bio             *string `json:"bio"`
	Address         *string `json:"address"`
	BirthDate       *string `json:"birth_date"`
}
