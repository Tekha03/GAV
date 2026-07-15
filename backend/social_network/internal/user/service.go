package user

import (
	"context"

	"social_network/internal/dog"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, email, passwordHash string) (*User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)

	FindDogsNearby(ctx context.Context, id uuid.UUID, centerLat, centerLon float64, radiusMeters float64) ([]*dog.Dog, error)

	UpdateLocation(ctx context.Context, userID uuid.UUID, input UpdateLocationInput) error
	SetLocationVisibility(ctx context.Context, userID uuid.UUID, input SetLocationVisibilityInput) error
}
