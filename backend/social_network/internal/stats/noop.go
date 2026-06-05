package stats

import (
	"context"

	"github.com/google/uuid"
)

type noopService struct{}

func NoopService() StatsService {
	return noopService{}
}

func ServiceOrNoop(services ...StatsService) StatsService {
	if len(services) > 0 && services[0] != nil {
		return services[0]
	}
	return NoopService()
}

func (noopService) UserStats(context.Context, uuid.UUID) (*UserStats, error) {
	return nil, ErrStatsNotFound
}

func (noopService) ProfileStats(context.Context, uuid.UUID) (*ProfileStats, error) {
	return nil, ErrStatsNotFound
}

func (noopService) IncrementPosts(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) IncrementFollowers(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) IncrementDogs(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) IncrementFollowings(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) DecrementPosts(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) DecrementFollowers(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) DecrementDogs(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) DecrementFollowings(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) PostStats(context.Context, uuid.UUID) (*PostStats, error) {
	return nil, ErrStatsNotFound
}

func (noopService) IncrementPostLikes(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) IncrementPostComments(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) DecrementPostLikes(context.Context, uuid.UUID) error {
	return nil
}

func (noopService) DecrementPostComments(context.Context, uuid.UUID) error {
	return nil
}
