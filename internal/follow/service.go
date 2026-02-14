package follow

import (
	"context"
)

type FollowService interface {
	Follow(ctx context.Context, follow Follow) error
	Unfollow(ctx context.Context, follow Follow) error
}
