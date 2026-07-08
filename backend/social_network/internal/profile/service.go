package profile

import (
	"context"

	"github.com/google/uuid"
)

type ProfileService interface {
	Create(ctx context.Context, userID uuid.UUID, input CreateProfileInput) (*UserProfile, error)
	GetByID(ctx context.Context, profileID uuid.UUID) (*UserProfile, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	Search(ctx context.Context, query string, limit int) ([]*UserProfile, error)
	Update(ctx context.Context, profileID uuid.UUID, input UpdateProfileInput) error
	Delete(ctx context.Context, profileID uuid.UUID) error
}
