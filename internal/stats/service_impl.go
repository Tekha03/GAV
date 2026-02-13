package stats

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrStatsNotFound = errors.New("user stats not found")

type Service struct {
	repo StatsRepository
}

func NewService(repo StatsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	stats, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return nil, ErrStatsNotFound
	}

	return stats, nil
}

func (s *Service) IncrementPosts(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementPosts(ctx, userID)
}

func (s *Service) IncrementFollowers(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementFollowers(ctx, userID)
}

func (s *Service) IncrementDogs(ctx context.Context, userID uuid.UUID) error {
	return s.repo.IncrementDogs(ctx, userID)
}

func (s *Service) DecrementPosts(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementPosts(ctx, userID)
}

func (s *Service) DecrementFollowers(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementFollowers(ctx, userID)
}

func (s *Service) DecrementDogs(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DecrementDogs(ctx, userID)
}
