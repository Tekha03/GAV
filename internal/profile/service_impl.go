package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrProfileAlreadyExists = errors.New("profile already exists")
	ErrProfileNotFound		= errors.New("profile not found")
	ErrInvalidUserID	   = errors.New("invalid user ID")
	ErrInvalidProfileID	  	= errors.New("invalid profile ID")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, input CreateProfileInput) (*UserProfile, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	profile := &UserProfile{
		UserID: userID,
		Name: input.Name,
		Surname: input.Surname,
		Username: input.Username,
		ProfilePhoto: input.ProfilePhoto,
		Bio: input.Bio,
		Address: input.Address,
		BirthDate: input.BirthDate,
	}

	if err := s.repo.Create(ctx, profile); err != nil {
		return nil, ErrProfileAlreadyExists
	}

	return profile, nil
}

func (s *Service) GetByID(ctx context.Context, profileID uuid.UUID) (*UserProfile, error) {
	if profileID == uuid.Nil {
		return nil, ErrInvalidProfileID
	}

	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, ErrProfileNotFound
	}

	return profile, nil
}

func (s *Service) Update(ctx context.Context, profileID uuid.UUID, input UpdateProfileInput) error {
	if profileID == uuid.Nil {
		return ErrInvalidProfileID
	}

	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return ErrProfileNotFound
	}

	if input.Name != nil {
		profile.Name = *input.Name
	}
	if input.Surname != nil {
		profile.Surname = *input.Surname
	}
	if input.Username != nil {
		profile.Username = *input.Username
	}
	if input.ProfilePhoto != nil {
		profile.ProfilePhoto = *input.ProfilePhoto
	}
	if input.Bio != nil {
		profile.Bio = *input.Bio
	}
	if input.Address != nil {
		profile.Address = *input.Address
	}
	if input.BirthDate != nil {
		profile.BirthDate = *input.BirthDate
	}

	return s.repo.Update(ctx, profile)
}

func (s *Service) Delete(ctx context.Context, profileID uuid.UUID) error {
	if profileID == uuid.Nil {
		return ErrInvalidProfileID
	}

	if err := s.repo.Delete(ctx, profileID); err != nil {
		return ErrProfileNotFound
	}

	return nil
}
