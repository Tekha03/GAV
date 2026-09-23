package postgres

import (
	"context"
	"social_network/internal/notification"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	*BaseRepository
}

func NewNotificationRepository(db *gorm.DB) (notification.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &NotificationRepository{BaseRepository: repo}, nil
}

func (r *NotificationRepository) Create(ctx context.Context, n *notification.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*notification.Notification, error) {
	var notifications []*notification.Notification
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&notification.Notification{}).
		Where("user_id = ? AND id = ?", userID, notificationID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) DeleteOld(ctx context.Context, userID uuid.UUID, ttl time.Duration) error {
	cutoff := time.Now().Add(-ttl)
	return r.db.WithContext(ctx).Where("user_id = ? AND created_at < ?", userID, cutoff).Delete(&notification.Notification{}).Error
}
