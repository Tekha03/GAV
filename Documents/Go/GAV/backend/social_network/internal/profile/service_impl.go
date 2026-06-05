package profile

import (
	"context"
	"strings"
	"unicode"

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

	username := normalizeUsername(input.Username)
	if !isValidUsername(username) {
		return nil, ErrInvalidUsername
	}

	profile := &UserProfile{
		UserID:          userID,
		Name:            input.Name,
		Surname:         input.Surname,
		Username:        username,
		ProfilePhotoUrl: input.ProfilePhotoUrl,
		Bio:             input.Bio,
		Address:         input.Address,
		BirthDate:       input.BirthDate,
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

func (s *service) Search(ctx context.Context, query string, limit int) ([]*UserProfile, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []*UserProfile{}, nil
	}

	if limit <= 0 || limit > 25 {
		limit = 10
	}

	return s.repo.Search(ctx, query, limit)
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
		username := normalizeUsername(*input.Username)
		if !isValidUsername(username) {
			return ErrInvalidUsername
		}
		profile.Username = username
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

	if err := s.repo.Update(ctx, profile); err != nil {
		return ErrProfileAlreadyExists
	}

	return nil
}

func normalizeUsername(username string) string {
	return strings.Trim(strings.ToLower(strings.TrimSpace(username)), "@")
}

func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	for _, r := range username {
		if r > unicode.MaxASCII {
			return false
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' {
			continue
		}
		return false
	}

	return true
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
