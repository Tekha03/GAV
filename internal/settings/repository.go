package settings

import "context"

type Repository interface {
	GetByUserID(ctx context.Context, userID uint) (*UserSettings, error)
	Create(ctx context.Context, settings *UserSettings) error
	Update(ctx context.Context, settings *UserSettings) error
}
