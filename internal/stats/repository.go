package stats

import (
	"context"

	"github.com/google/uuid"
)

type StatsRepository interface {
    Create(ctx context.Context, st *UserStats) error
    Delete(ctx context.Context, userID uuid.UUID) error
    GetByUserID(ctx context.Context, userID uuid.UUID) (*UserStats, error)

    IncrementPosts(ctx context.Context, userID uuid.UUID) error
    IncrementDogs(ctx context.Context, userID uuid.UUID) error
    IncrementFollowers(ctx context.Context, userID uuid.UUID) error
    IncrementFollowings(ctx context.Context, userID uuid.UUID) error

    DecrementPosts(ctx context.Context, userID uuid.UUID) error
    DecrementFollowers(ctx context.Context, userID uuid.UUID) error
    DecrementDogs(ctx context.Context, userID uuid.UUID) error
    DecrementFollowings(ctx context.Context, userID uuid.UUID) error
}
