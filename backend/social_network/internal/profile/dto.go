package profile

type CreateProfileInput struct {
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Username        string `json:"username"`
	ProfilePhotoUrl string `json:"profile_photo_url"`
	Bio             string `json:"bio"`
	Address         string `json:"address"`
	BirthDate       string `json:"birth_date"`
}

type UpdateProfileInput struct {
	Name            *string `json:"name"`
	Surname         *string `json:"surname"`
	Username        *string `json:"username"`
	ProfilePhotoUrl *string `json:"profile_photo_url"`
	Bio             *string `json:"bio"`
	Address         *string `json:"address"`
	BirthDate       *string `json:"birth_date"`
}
