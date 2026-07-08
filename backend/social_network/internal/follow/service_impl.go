package follow

import (
	"context"
	"social_network/internal/stats"

	"github.com/google/uuid"
)

type service struct {
	repo        Repository
	statService stats.StatsService
}

func NewService(repo Repository, statService ...stats.StatsService) (FollowService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo, statService: stats.ServiceOrNoop(statService...)}, nil
}

func (s *service) Follow(ctx context.Context, follow Follow) error {
	if follow.FollowerID == follow.FollowingID {
		return ErrCannotFollowYourself
	}

	alreadyFollowing, err := s.repo.FollowerExists(ctx, follow)
	if err != nil {
		return err
	}

	if alreadyFollowing {
		return ErrAlreadyFollowing
	}

	if err = s.statService.IncrementFollowings(ctx, follow.FollowerID); err != nil {
		return err
	}

	if err = s.statService.IncrementFollowers(ctx, follow.FollowingID); err != nil {
		return err
	}

	return s.repo.Follow(ctx, follow)
}

func (s *service) Unfollow(ctx context.Context, follow Follow) error {
	if err := s.statService.DecrementFollowings(ctx, follow.FollowerID); err != nil {
		return err
	}

	if err := s.statService.DecrementFollowers(ctx, follow.FollowingID); err != nil {
		return err
	}
	return s.repo.Unfollow(ctx, follow)
}

func (s *service) GetFollowers(ctx context.Context, userID uuid.UUID) ([]Follow, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	followers, err := s.repo.GetFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (s *service) GetFollowing(ctx context.Context, userID uuid.UUID) ([]Follow, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	following, err := s.repo.GetFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	return following, nil
}
