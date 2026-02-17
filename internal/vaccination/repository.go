package vaccination

import (
	"context"

	"github.com/google/uuid"
)

type VaccinationRepository interface {
	Create(ctx context.Context, v *Vaccination) error
	Update(ctx context.Context, v *Vaccination) error
	// AddVaccination(ctx context.Context, v *Vaccination) error
	Delete(ctx context.Context, ID uuid.UUID) error
	GetByDogID(ctx context.Context, dogID uuid.UUID) ([]Vaccination, error)
	GetByID(ctx context.Context, ID uuid.UUID) (*Vaccination, error)
}
