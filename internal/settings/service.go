package settings

import "context"

type SettingsService interface {
    Get(ctx context.Context, userID uint) (*UserSettings, error)
    Update(ctx context.Context, userID uint, input UpdateSettingsInput) error
}
