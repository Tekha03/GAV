package vaccination

import (
	"context"

	"github.com/google/uuid"
)

type VaccinationService interface {
	Create(ctx context.Context, dogID uuid.UUID, input CreateVaccinationInput) (*Vaccination, error)
	// AddVaccination(ctx context.Context, ownerID, dogID uuid.UUID, v Vaccination) (*Vaccination, error)
	ListByDogID(ctx context.Context, dogID uuid.UUID) ([]*Vaccination, error)
	Update(ctx context.Context, vaccinationID, dogID uuid.UUID, input UpdateVaccinationInput) error
	Delete(ctx context.Context, vaccinationID uuid.UUID) error
}
