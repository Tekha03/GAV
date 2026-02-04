package dog

import (
	"context"
	"errors"
)

var (
	ErrDogAccessDenied = errors.New("dog access denied")
)

type DogService struct {
	repo 	DogRepository
}

func NewDogService(repo DogRepository) *DogService {
	return &DogService{repo: repo}
}

func (s *DogService) Create(ctx context.Context, ownerID uint, input CreateDogInput) (*Dog, error) {
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

func (s *DogService) Update(ctx context.Context, ownerID, dogID uint, input UpdateDogInput) error {

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

func (s *DogService) GetPublic(ctx context.Context, dogID uint) (*Dog, error) {
	return s.repo.GetByID(ctx, dogID)
}

func (s *DogService) GetPrivate(ctx context.Context, dogID, ownerID uint) (*Dog, error) {
	dog, err := s.repo.GetByID(ctx, dogID)
	if err != nil {
		return nil, err
	}

	if dog.OwnerID != ownerID {
		return nil, ErrDogAccessDenied
	}

	return dog, nil
}