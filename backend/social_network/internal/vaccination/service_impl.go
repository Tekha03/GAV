package vaccination

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) (VaccinationService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo}, nil
}

func (s *service) Create(ctx context.Context, dogID uuid.UUID, input CreateVaccinationInput) (*Vaccination, error) {
	vaccination := Vaccination{
		DogID:     dogID,
		Name:      input.Name,
		DoneAt:    input.DoneAt,
		NextDueAt: input.NextDueAt,
		Notes:     input.Notes,
	}

	if err := s.repo.Create(ctx, &vaccination); err != nil {
		return nil, err
	}

	return &vaccination, nil
}

func (s *service) ListByDogID(ctx context.Context, dogID uuid.UUID) ([]*Vaccination, error) {
	if dogID == uuid.Nil {
		return nil, ErrDogIDEmpty
	}

	return s.repo.ListByDogID(ctx, dogID)
}

func (s *service) Update(ctx context.Context, vaccinationID, dogID uuid.UUID, input UpdateVaccinationInput) error {
	current_vaccine, err := s.repo.GetByID(ctx, vaccinationID)
	if err != nil {
		return err
	}

	if current_vaccine.DogID != dogID {
		return ErrVaccAccessDenied
	}

	if input.Name != nil {
		current_vaccine.Name = *input.Name
	}

	if input.DoneAt != nil {
		current_vaccine.DoneAt = *input.DoneAt
	}

	if input.NextDueAt != nil {
		current_vaccine.NextDueAt = input.NextDueAt
	}

	if input.Notes != nil {
		current_vaccine.Notes = *input.Notes
	}

	return s.repo.Update(ctx, current_vaccine)
}

func (s *service) Delete(ctx context.Context, vaccinationID uuid.UUID) error {
	vaccination, err := s.repo.GetByID(ctx, vaccinationID)
	if err != nil {
		return err
	}

	if vaccination.ID != vaccinationID {
		return ErrVaccAccessDenied
	}

	return s.repo.Delete(ctx, vaccinationID)
}
