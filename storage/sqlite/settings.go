package sqlite

import (
	"context"
	"errors"
	"fmt"
	"gav/internal/settings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingsRepository struct {
	*BaseRepository
}

func NewSettingsRepository(db *gorm.DB) settings.Repository {
	return &SettingsRepository{NewBaseRepository(db)}
}

func (r *SettingsRepository) Create(ctx context.Context, settins *settings.UserSettings) error {
	if err := r.DB(ctx).Create(settins).Error; err != nil {
		return fmt.Errorf("settings repository: create: %w", err)
	}

	return  nil
}

func (r *SettingsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*settings.UserSettings, error) {
	var userUsettings settings.UserSettings

	if err := r.DB(ctx).First(&userUsettings, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSettingsNotFound
		}

		return nil, fmt.Errorf("settings repository: get by user id: %w", err)
	}

	return &userUsettings, nil
}

func (r *SettingsRepository) Update(ctx context.Context, userSettings *settings.UserSettings) error {
	result := r.DB(ctx).
		Model(&settings.UserSettings{}).
		Where("user_id = ?", userSettings.UserID).
		Updates(userSettings)

	if result.Error != nil {
		return fmt.Errorf("settings repository: update: %w", result.Error)
	}

	return nil
}

