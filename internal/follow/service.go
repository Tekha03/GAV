package follow

import (
	"context"
	"errors"
)

type FollowService struct {
	repo FollowRepository
}

func (fs *FollowService) Follow(ctx context.Context, follower, following uint) error {
	if follower == following {
		return errors.New("you cannot follow yourself")
	}

	return fs.repo.Follow(ctx, follower, following)
}
