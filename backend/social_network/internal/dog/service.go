package dog

import (
    "context"
    "github.com/google/uuid"
)

type DogService interface {
	Create(ctx context.Context, ownerID uuid.UUID, input CreateDogInput) (*Dog, error)
    Update(ctx context.Context, ownerID, dogID uuid.UUID, input UpdateDogInput) error
    Delete(ctx context.Context, ownerID, dogID uuid.UUID) error

    GetPublic(ctx context.Context, dogID uuid.UUID) (*Dog, error)
    GetPrivate(ctx context.Context, ownerID, dogID uuid.UUID) (*Dog, error)

	// later for analytics
	// GetStatusHistory(ownerID uint, dogID uint) ([]StatusChange, error)
}
