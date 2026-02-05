package vaccination

import "context"

type Service interface {
	Create(ctx context.Context, dogID uint, v CreateVaccinationInput) error
	// AddVaccination(ctx context.Context, ownerID, dogID uint, v Vaccination) (*Vaccination, error)
	Update(ctx context.Context, ID, dogID uint, v UpdateVaccinationInput) error
	Delete(ctx context.Context, ID, dogID uint, vaccinationID uint) error
	GetByDogID(ctx context.Context, dogID uint) ([]Vaccination, error)
}