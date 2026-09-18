package postgres

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
	return tokens, r.DB(ctx).Distinct().Where("user_id = ?", userID).Find(&tokens).Error
}

func (r *DeviceTokenRepository) SaveForUser(ctx context.Context, userID uuid.UUID, token string) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtextextended(?, 0))", token).Error; err != nil {
			return err
		}
		if err := tx.Where("token = ?", token).Delete(&device.DeviceToken{}).Error; err != nil {
			return err
		}
		return tx.Create(&device.DeviceToken{UserID: userID, Token: token}).Error
	})
}

func (r *DeviceTokenRepository) DeleteForUser(ctx context.Context, userID uuid.UUID, token string) error {
	return r.DB(ctx).Where("user_id = ? AND token = ?", userID, token).Delete(&device.DeviceToken{}).Error
}
