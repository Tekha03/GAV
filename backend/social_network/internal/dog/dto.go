package dog

type CreateDogInput struct {
	Name     string `json:"name"`
	Breed    string `json:"breed"`
	Gender   Gender `json:"gender"`
	Age      Age    `json:"age"`
	Status   Status `json:"status"`
	PhotoUrl string `json:"photo_url"`
	Notes    string `json:"notes"`
}

type UpdateDogInput struct {
	Name     *string `json:"name"`
	Breed    *string `json:"breed"`
	Gender   *Gender `json:"gender"`
	Age      *Age    `json:"age"`
	Status   *Status `json:"status"`
	PhotoUrl *string `json:"photo_url"`
	Notes    *string `json:"notes"`
}
