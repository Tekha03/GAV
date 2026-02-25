package feed

import (
	"context"
	"gav/internal/post"
	"time"

	"github.com/google/uuid"
)

type service struct {
	postRepo post.Repository
}

func NewService(postRepo post.Repository) (FeedService, error) {
	if postRepo == nil {
		return nil, post.ErrRepoNil
	}

	return &service{postRepo: postRepo}, nil
}

func (s *service) GetFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*post.Post, time.Time, error) {
	limit = max(limit, 20)
	limit = min(limit, 100)

	posts, err := s.postRepo.ListFeed(ctx, userID, before, limit + 1)
	if err != nil {
		return nil, time.Time{}, err
	}

	var nextTime time.Time
	if len(posts) > limit {
		nextTime = posts[limit].CreatedAt
		posts = posts[:limit]
	}

	return posts, nextTime, nil
}
