package stats

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrStatsNotFound = errors.New("user stats not found")

type service struct {
	repo Repository
}

func NewService(repo Repository) StatsService {
	return &service{repo: repo}
}

func (s *service) UserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	return s.repo.GetUserStats(ctx, userID)
}

func (s *service) ProfileStats(ctx context.Context, userID uuid.UUID) (*ProfileStats, error) {
	userStats, err := s.repo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, ErrStatsNotFound
	}

	return &ProfileStats{
		UserID: userStats.UserID,
		PostCount: userStats.PostCount,
		Followers: userStats.Followers,
		Followings: userStats.Followings,
	}, nil
}

func (s *service) IncrementPosts(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementPosts(ctx, userID)
}

func (s *service) IncrementFollowers(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementFollowers(ctx, userID)
}

func (s *service) IncrementDogs(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementDogs(ctx, userID)
}

func (s *service) IncrementFollowings(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementFollowings(ctx, userID)
}

func (s *service) DecrementPosts(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementPosts(ctx, userID)
}

func (s *service) DecrementFollowers(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementFollowers(ctx, userID)
}

func (s *service) DecrementDogs(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementDogs(ctx, userID)
}

func (s *service) DecrementFollowings(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementFollowings(ctx, userID)
}

func (s *service) PostStats(ctx context.Context, postID uuid.UUID) (*PostStats, error) {
	return s.repo.GetPostStats(ctx, postID)
}

func (s *service) IncrementPostLikes(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementPostLikes(ctx, userID)
}

func (s *service) DecrementPostLikes(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementPostLikes(ctx, userID)
}

func (s *service) IncrementPostComments(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementPostComments(ctx, userID)
}

func (s *service) DecrementPostComments(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementPostComments(ctx, userID)
}
