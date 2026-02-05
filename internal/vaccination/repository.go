package vaccination

import "context"

type VaccinationRepository interface {
	Create(ctx context.Context, v *Vaccination) (*Vaccination, error)
	Update(ctx context.Context, v *Vaccination) error
	// AddVaccination(ctx context.Context, v *Vaccination) error
	Delete(ctx context.Context, ID uint) error
	GetByDogID(ctx context.Context, dogID uint) ([]Vaccination, error)
	GetByID(ctx context.Context, ID uint) (*Vaccination, error)
}