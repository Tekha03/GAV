package dog

import "github.com/google/uuid"

type Dog struct {
	ID      uuid.UUID `json:"id"`
	OwnerID uuid.UUID `json:"owner_id"`

	Name     string `json:"name"`
	Breed    string `json:"breed"`
	PhotoUrl string `json:"photo_url"`
	Notes    string `json:"notes"`

	Status Status `json:"status"`
	Age    Age    `json:"age"`
	Gender Gender `json:"gender"`

	Lat *float64 `json:"lat,omitempty"`
	Lon *float64 `json:"lon,omitempty"`
}

func NewDog(
	ownerID uuid.UUID,
	name string,
	breed string,
	gender Gender,
	status Status,
	age Age,
	photoURL string,
	notes string,
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
		ID:       uuid.New(),
		OwnerID:  ownerID,
		Name:     name,
		Breed:    breed,
		Gender:   gender,
		Status:   status,
		Age:      age,
		PhotoUrl: photoURL,
		Notes:    notes,
	}, nil
}
