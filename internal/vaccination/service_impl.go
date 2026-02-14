package vaccination

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrVaccAccessDenied = errors.New("vaccination access denied")
)

type VaccinationService struct {
	repo VaccinationRepository
}

func NewService(repo VaccinationRepository) *VaccinationService {
	return &VaccinationService{
		repo: repo,
	}
}

func (s *VaccinationService) Create(ctx context.Context, dogID uuid.UUID, vi *CreateVaccinationInput)(*Vaccination, error ){
	vaccination := Vaccination{
		DogID: dogID,
		Name: vi.Name,
		DoneAt: vi.DoneAt,
		NextDueAt: vi.NextDueAt,
		Notes: vi.Notes,
	}

	if err := s.repo.Create(ctx, &vaccination); err != nil {
		return nil, err
	}

	return &vaccination, nil
}

func (s *VaccinationService) Update(ctx context.Context, ID, dogID uuid.UUID, input UpdateVaccinationInput) error {
	vac, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	if vac.DogID != dogID {
		return ErrVaccAccessDenied
	}

	if input.Name != nil {
		vac.Name = *input.Name
	}

	if input.DoneAt != nil {
		vac.DoneAt = *input.DoneAt
	}

	if input.NextDueAt != nil {
		vac.NextDueAt = input.NextDueAt
	}

	if input.Notes != nil {
		vac.Notes = *input.Notes
	}

	return nil
}
