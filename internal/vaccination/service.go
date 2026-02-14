package vaccination

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, dogID uuid.UUID, v CreateVaccinationInput) error
	// AddVaccination(ctx context.Context, ownerID, dogID uuid.UUID, v Vaccination) (*Vaccination, error)
	Update(ctx context.Context, ID, dogID uuid.UUID, v UpdateVaccinationInput) error
	Delete(ctx context.Context, ID, dogID uuid.UUID, vaccinationID uuid.UUID) error
	GetByDogID(ctx context.Context, dogID uuid.UUID) ([]Vaccination, error)
}
