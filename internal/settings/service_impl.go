package settings

import (
	"context"
	"errors"
)

var (
	ErrSettingsNotFound = errors.New("settings not found")
	ErrInvalidUserID	= errors.New("invalid user ID")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, userID uint) (*UserSettings, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	settings, err := s.repo.GetByUserID(ctx, userID)
	if err == nil {
		return settings, nil
	}

	defaultSettings := &UserSettings{
		UserID:			userID,
		ProfilePrivacy:	 false,
		ShowLocation:	true,
		AllowMessages:	true,
	}

	if err := s.repo.Create(ctx, defaultSettings); err != nil {
		return nil, err
	}

	return defaultSettings, nil
}

func (s *Service) Update(ctx context.Context, userID uint, input UpdateSettingsInput) error {
	if userID == 0 {
		return ErrInvalidUserID
	}

	settings, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return ErrSettingsNotFound
	}

	if input.ProfilePrivacy != nil {
		settings.ProfilePrivacy = *input.ProfilePrivacy
	}
	if input.ShowLocation != nil {
		settings.ShowLocation = *input.ShowLocation
	}
	if input.AllowMessages != nil {
		settings.AllowMessages = *input.AllowMessages
	}

	return s.repo.Update(ctx, settings)
}
