package dog

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDogAccessDenied = errors.New("dog access denied")
)

type service struct {
	repo DogRepository
}

func NewDogService(repo DogRepository) DogService {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, ownerID uuid.UUID, input CreateDogInput) (*Dog, error) {
	dog := NewDog(
		ownerID,
		input.Name,
		input.Breed,
		input.Gender,
		input.Status,
		input.Age,
		input.PhotoUrl,
	)

	if err := s.repo.Create(ctx, dog); err != nil {
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

	return s.repo.Delete(ctx, dogID)
}

func (s *service) UpdateLocation(ctx context.Context, ownerID, dogID uuid.UUID, lat, lon float64) error {
	dog, err := s.repo.GetByID(ctx, dogID)
	if err != nil {
		return err
	}

	if dog.OwnerID != ownerID {
		return ErrDogAccessDenied
	}

	dog.Lat = &lat
	dog.Lon = &lon

	return s.repo.Update(ctx, dog)
}

func (s *service) SetLocationVisibility(ctx context.Context, ownerID, dogID uuid.UUID, visible bool) error {
	dog, err := s.repo.GetByID(ctx, dogID)
	if err != nil {
		return err
	}

	if dog.OwnerID != ownerID {
		return ErrDogAccessDenied
	}

	dog.LocationVisible = visible

	return s.repo.Update(ctx, dog)
}

func (s *service) GetPublic(ctx context.Context, dogID uuid.UUID) (*Dog, error) {
	return s.repo.GetByID(ctx, dogID)
}

func (s *service) GetPrivate(ctx context.Context, dogID, ownerID uuid.UUID) (*Dog, error) {
	dog, err := s.repo.GetByID(ctx, dogID)
	if err != nil {
		return nil, err
	}

	if dog.OwnerID != ownerID {
		return nil, ErrDogAccessDenied
	}

	return dog, nil
}
