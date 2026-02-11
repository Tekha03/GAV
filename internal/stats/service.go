package stats

import "context"

type StatsService interface {
    Get(ctx context.Context, userID uint) (*UserStats, error)

    IncrementPosts(ctx context.Context, userID uint) error
    IncrementFollowers(ctx context.Context, userID uint) error
    IncrementDogs(ctx context.Context, userID uint) error

    DecrementPosts(ctx context.Context, userID uint) error
    DecrementFollowers(ctx context.Context, userID uint) error
    DecrementDogs(ctx context.Context, userID uint) error
}
