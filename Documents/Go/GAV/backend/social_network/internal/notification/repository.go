package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, n *Notification) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Notification, error)
	MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error
	DeleteOld(ctx context.Context, userID uuid.UUID, ttl time.Duration) error
}