package profile

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) (ProfileService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo}, nil
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, input CreateProfileInput) (*UserProfile, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	profile := &UserProfile{
		UserID: userID,
		Name: input.Name,
		Surname: input.Surname,
		Username: input.Username,
		ProfilePhotoUrl: input.ProfilePhotoUrl,
		Bio: input.Bio,
		Address: input.Address,
		BirthDate: input.BirthDate,
	}

	if err := s.repo.Create(ctx, profile); err != nil {
		return nil, ErrProfileAlreadyExists
	}

	return profile, nil
}

func (s *service) GetByID(ctx context.Context, profileID uuid.UUID) (*UserProfile, error) {
	if profileID == uuid.Nil {
		return nil, ErrInvalidProfileID
	}

	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, ErrProfileNotFound
	}

	return profile, nil
}

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	profile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *service) Update(ctx context.Context, profileID uuid.UUID, input UpdateProfileInput) error {
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
	if input.ProfilePhotoUrl != nil {
		profile.ProfilePhotoUrl = *input.ProfilePhotoUrl
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

func (s *service) Delete(ctx context.Context, profileID uuid.UUID) error {
	if profileID == uuid.Nil {
		return ErrInvalidProfileID
	}

	if err := s.repo.Delete(ctx, profileID); err != nil {
		return ErrProfileNotFound
	}

	return nil
}
