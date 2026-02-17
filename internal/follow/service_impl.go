package follow

import (
	"context"
	"errors"
)

var (
	ErrCannotFollowYourself = errors.New("you cannot follow yourself.")
	ErrAlreadyFollowing = errors.New("already following")
)

type service struct {
	repo Repository
}

func NewService(repo Repository) FollowService {
	return &service{repo: repo}
}

func (s *service) Follow(ctx context.Context, follow Follow) error {
	if follow.FollowerID == follow.FollowingID {
		return ErrCannotFollowYourself
	}

	alreadyFollowing, err := s.repo.FollowerExists(ctx, follow)
	if err != nil {
		return err
	}

	if  alreadyFollowing {
		return ErrAlreadyFollowing
	}

	return s.repo.Follow(ctx, follow)
}

func (s *service) Unfollow(ctx context.Context, follow Follow) error {
	return s.repo.Unfollow(ctx, follow)
}
