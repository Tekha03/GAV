package follow

import (
	"context"

	"github.com/google/uuid"
)

type FollowService interface {
	Follow(ctx context.Context, follow Follow) error
	Unfollow(ctx context.Context, follow Follow) error
	GetFollowers(ctx context.Context, userID uuid.UUID) ([]Follow, error)
	GetFollowing(ctx context.Context, userID uuid.UUID) ([]Follow, error)
}
