package feed

import (
	"context"
	"social_network/internal/post"
	"time"

	"github.com/google/uuid"
)

type FeedService interface {
	GetFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*post.Post, time.Time, error)
}
