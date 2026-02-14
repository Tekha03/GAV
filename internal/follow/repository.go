package follow

import "context"

type Repository interface {
	Follow(ctx context.Context, follow Follow) error
	Unfollow(ctx context.Context, follow Follow) error
	FollowerExists(ctx context.Context, follow Follow) (bool, error)
}
