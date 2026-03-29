package settings

import (
	"context"

	"github.com/google/uuid"
)

type SettingsService interface {
    Get(ctx context.Context, userID uuid.UUID) (*UserSettings, error)
    Update(ctx context.Context, userID uuid.UUID, input UpdateSettingsInput) error
}
