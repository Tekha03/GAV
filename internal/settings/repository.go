package settings

import "context"

type Repository interface {
	Create(ctx context.Context, settings *UserSettings) error
	GetByUserID(ctx context.Context, userID uint) (*UserSettings, error)
	Update(ctx context.Context, settings *UserSettings) error
}
