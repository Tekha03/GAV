package dog

import (
	"context"
	"social_network/internal/stats"

	"github.com/google/uuid"
)

type service struct {
	repo        Repository
	statService stats.StatsService
}

func NewService(repo Repository, statService ...stats.StatsService) (DogService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo, statService: stats.ServiceOrNoop(statService...)}, nil
}

func (s *service) Create(ctx context.Context, ownerID uuid.UUID, input CreateDogInput) (*Dog, error) {
	dog, err := NewDog(
		ownerID,
		input.Name,
		input.Breed,
		input.Gender,
		input.Status,
		input.Age,
		input.PhotoUrl,
		input.Notes,
	)
	if err != nil {
		return nil, err
	}

	if err = s.repo.Create(ctx, dog); err != nil {
		return nil, err
	}

	if err = s.statService.IncrementDogs(ctx, ownerID); err != nil {
		return nil, err
	}

	return dog, nil
}

func (s *service) Update(ctx context.Context, ownerID, dogID uuid.UUID, input UpdateDogInput) error {

	dog, err := s.repo.GetByID(ctx, dogID)
	if err != nil {
		return err
	}

	if dog.OwnerID != ownerID {
		return ErrDogAccessDenied
	}

	if input.Name != nil {
		dog.Name = *input.Name
	}

	if input.Breed != nil {
		dog.Breed = *input.Breed
	}

	if input.PhotoUrl != nil {
		dog.PhotoUrl = *input.PhotoUrl
	}

	if input.Age != nil {
		dog.Age = *input.Age
	}

	if input.Gender != nil {
		dog.Gender = *input.Gender
	}

	if input.Status != nil {
		dog.Status = *input.Status
	}

	if input.Notes != nil {
		dog.Notes = *input.Notes
	}

	return s.repo.Update(ctx, dog)
}

func (s *service) Delete(ctx context.Context, ownerID, dogID uuid.UUID) error {
	dog, err := s.repo.GetByID(ctx, dogID)
	if err != nil {
		return err
	}

	if dog.OwnerID != ownerID {
		return ErrDogAccessDenied
	}

	if err = s.statService.DecrementDogs(ctx, ownerID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, dogID)
}

func (s *service) GetPublic(ctx context.Context, dogID uuid.UUID) (*Dog, error) {
	return s.repo.GetByID(ctx, dogID)
}

func (s *service) GetPrivate(ctx context.Context, ownerID, dogID uuid.UUID) (*Dog, error) {
	dog, err := s.repo.GetByID(ctx, dogID)
	if err != nil {
		return nil, err
	}

	if dog.OwnerID != ownerID {
		return nil, ErrDogAccessDenied
	}

	return dog, nil
}

func (s *service) ListByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*Dog, error) {
	if ownerID == uuid.Nil {
		return nil, ErrOwnerIDNil
	}

	return s.repo.GetByOwnerID(ctx, ownerID)
}
