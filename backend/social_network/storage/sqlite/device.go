package sqlite

import (
	"context"
	"social_network/internal/device"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceTokenRepository struct {
	*BaseRepository
}

func NewDeviceRepo(db *gorm.DB) (device.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &DeviceTokenRepository{BaseRepository: repo}, nil
}

func (r *DeviceTokenRepository) Create(ctx context.Context, t *device.DeviceToken) error {
	return r.db.Create(t).Error
}

func (r *DeviceTokenRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*device.DeviceToken, error) {
	var tokens []*device.DeviceToken
	return tokens, r.db.Where("user_id = ?", userID).Find(&tokens).Error
}
