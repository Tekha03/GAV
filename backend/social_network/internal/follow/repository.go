package follow

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Follow(ctx context.Context, follow Follow) error
	Unfollow(ctx context.Context, follow Follow) error
	FollowerExists(ctx context.Context, follow Follow) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID) ([]Follow, error)
	GetFollowing(ctx context.Context, userID uuid.UUID) ([]Follow, error)
}
