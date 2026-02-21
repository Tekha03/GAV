package stats

import (
	"context"

	"github.com/google/uuid"
)

type StatsService interface {
    UserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
    ProfileStats(ctx context.Context, userID uuid.UUID) (*ProfileStats, error)

    IncrementPosts(ctx context.Context, userID uuid.UUID) error
    IncrementFollowers(ctx context.Context, userID uuid.UUID) error
    IncrementDogs(ctx context.Context, userID uuid.UUID) error
    IncrementFollowings(ctx context.Context, userID uuid.UUID) error

    DecrementPosts(ctx context.Context, userID uuid.UUID) error
    DecrementFollowers(ctx context.Context, userID uuid.UUID) error
    DecrementDogs(ctx context.Context, userID uuid.UUID) error
    DecrementFollowings(ctx context.Context, userID uuid.UUID) error

    PostStats(ctx context.Context, postID uuid.UUID) (*PostStats, error)
    IncrementPostLikes(ctx context.Context, postID uuid.UUID) error
    IncrementPostComments(ctx context.Context, postID uuid.UUID) error
    DecrementPostLikes(ctx context.Context, postID uuid.UUID) error
    DecrementPostComments(ctx context.Context, postID uuid.UUID) error
}
