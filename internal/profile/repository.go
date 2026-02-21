package profile

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, profile *UserProfile) error
	GetByID(ctx context.Context, profileID uuid.UUID) (*UserProfile, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	Update(ctx context.Context, profile *UserProfile) error
	Delete(ctx context.Context, profileID uuid.UUID) error
}
