package vaccination

import "context"

type VaccinationService interface {
	Create(ctx context.Context, ownerID, dogID uint, v CreateVaccinationInput) error
	// AddVaccination(ctx context.Context, ownerID, dogID uint, v Vaccination) (*Vaccination, error)
	UpdateVaccination(ctx context.Context, ownerID, dogID uint, v UpdateVaccinationInput) error
	DeleteVaccination(ctx context.Context, ownerID, dogID uint, vaccinationID uint) error
	GetVaccinations(ctx context.Context, ownerID, dogID uint) ([]Vaccination, error)
}