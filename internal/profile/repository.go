package profile

import "context"

type Repository interface {
	Create(ctx context.Context, profile *UserProfile) error
	GetByID(ctx context.Context, profileID uint) (*UserProfile, error)
	Update(ctx context.Context, profile *UserProfile) error
	Delete(ctx context.Context, profileID uint) error
}
