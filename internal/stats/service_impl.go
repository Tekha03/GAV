package stats

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrStatsNotFound = errors.New("user stats not found")

type service struct {
	repo StatsRepository
}

func NewService(repo StatsRepository) StatsService {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	stats, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return nil, ErrStatsNotFound
	}

	return stats, nil
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

func (s *service) DecrementPosts(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementPosts(ctx, userID)
}

func (s *service) DecrementFollowers(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementFollowers(ctx, userID)
}

func (s *service) DecrementDogs(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementDogs(ctx, userID)
}
