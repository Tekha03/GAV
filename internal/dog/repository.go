package dog

import (
	"context"

	"github.com/google/uuid"
)

type DogRepository interface {
	Create(ctx context.Context, dog *Dog) error
	Update(ctx context.Context, dog *Dog) error
	Delete(ctx context.Context, ID uuid.UUID) error
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*Dog, error)
	GetByID(ctx context.Context, ID uuid.UUID) (*Dog, error)
}
