package memory

import (
	"context"
	"errors"
	"gav/internal/stats"
	"sync"
)

var (
	ErrStatExist = errors.New("stat exist in repository")
	ErrStatNotFound = errors.New("stat not found")
)

type StatsRepository struct {
	mu 		sync.RWMutex
	stats	map[uint]*stats.UserStats	
}

func NewStatsReposirory() *StatsRepository {
	return &StatsRepository{
		stats: make(map[uint]*stats.UserStats),
	}
}

func (s *StatsRepository) Create(ctx context.Context, st *stats.UserStats) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[st.UserID]; found {
		return ErrStatExist
	}

	s.stats[st.UserID] = st
	return nil
}

func (s *StatsRepository) GetByUserID(ctx context.Context, userID uint) (*stats.UserStats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return nil, ErrStatNotFound
	}

	return s.stats[userID], nil
}

func (s *StatsRepository) Delete(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	delete(s.stats, userID)
	return nil
}

func (s *StatsRepository) IncrementPosts(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	s.stats[userID].PostCount++

	return nil
}

func (s *StatsRepository) IncrementFollowers(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	s.stats[userID].Followers++

	return nil
}

func (s *StatsRepository) IncrementFolowings(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	s.stats[userID].Followings++

	return nil
}

func (s *StatsRepository) IncrementDogs(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	s.stats[userID].DogsCount++

	return nil
}

func (s *StatsRepository) DecrementDogs(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	s.stats[userID].DogsCount--

	return nil
}

func (s *StatsRepository) DecrementPosts(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	s.stats[userID].PostCount--

	return nil
}

func (s *StatsRepository) DecrementFollowers(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	s.stats[userID].Followers--

	return nil
}

func (s *StatsRepository) DecrementFollowings(ctx context.Context, userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.stats[userID]; !found {
		return ErrStatNotFound
	}

	s.stats[userID].Followings--

	return nil
}