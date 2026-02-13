package settings

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, settings *UserSettings) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*UserSettings, error)
	Update(ctx context.Context, settings *UserSettings) error
}
