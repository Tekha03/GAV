package user

import (
	"context"

	"social_network/internal/dog"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindWalkingNearby(ctx context.Context, centerLat, centerLon float64, radiusMeters float64) ([]*dog.Dog, error)
}
