package sqlite

import (
	"context"
	"errors"
	"fmt"
	"social_network/internal/profile"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	*BaseRepository
}

func NewProfileRepository(db *gorm.DB) (profile.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &ProfileRepository{BaseRepository: repo}, nil
}

func (r *ProfileRepository) Create(ctx context.Context, userProfile *profile.UserProfile) error {
	if err := r.DB(ctx).Create(userProfile).Error; err != nil {
		return fmt.Errorf("profile repository: create: %w", err)
	}

	return nil
}

func (r *ProfileRepository) GetByID(ctx context.Context, profileID uuid.UUID) (*profile.UserProfile, error) {
	var userProfile profile.UserProfile

	err := r.DB(ctx).First(&userProfile, "id = ?", profileID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, profile.ErrProfileNotFound
		}

		return nil, fmt.Errorf("profile repository: get by id: %w", err)
	}

	return &userProfile, nil
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*profile.UserProfile, error) {
	var userProfile profile.UserProfile

	err := r.DB(ctx).Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, profile.ErrProfileNotFound
		}

		return nil, fmt.Errorf("profile repository: get by user id: %w", err)
	}

	return &userProfile, nil
}

func (r *ProfileRepository) Update(ctx context.Context, userProfile *profile.UserProfile) error {
	err := r.DB(ctx).Save(userProfile).Error

	if err != nil {
		return fmt.Errorf("profile repository: update: %w", err)
	}

	return nil
}

func (r *ProfileRepository) Delete(ctx context.Context, profileID uuid.UUID) error {
	result := r.DB(ctx).Delete(&profile.UserProfile{}, "id = ?", profileID)

	if result.Error != nil {
		return fmt.Errorf("profile repository: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return profile.ErrProfileNotFound
	}

	return nil
}
