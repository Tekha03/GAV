package device

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *DeviceToken) error
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*DeviceToken, error)
}
