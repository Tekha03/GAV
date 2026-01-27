package follow

import "context"

type FollowRepository interface {
	Follow(ctx context.Context, follower, following uint) error
	Unfollow(ctx context.Context, follower, following uint) error
}
