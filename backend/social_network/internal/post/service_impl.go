package post

import (
	"context"
	"social_network/internal/stats"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo        Repository
	statService stats.StatsService
}

func NewService(repo Repository, statService ...stats.StatsService) (PostService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo, statService: stats.ServiceOrNoop(statService...)}, nil
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, content string, imageUrl string) (*Post, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}

	post := &Post{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		ImageUrl:  imageUrl,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	if err := s.statService.IncrementPosts(ctx, userID); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) GetByID(ctx context.Context, postID uuid.UUID) (*Post, error) {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (s *service) ListByUser(ctx context.Context, userID uuid.UUID) ([]*Post, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) GetFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*Post, time.Time, error) {
	posts, err := s.repo.ListFeed(ctx, userID, before, limit)
	if err != nil {
		return nil, time.Time{}, err
	}

	var nextCursor time.Time
	hasMore := len(posts) > limit

	if hasMore {
		nextCursor = posts[limit].CreatedAt
		posts = posts[:limit]
	}

	return posts, nextCursor, nil
}

func (s *service) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	if post == nil {
		return ErrPostNotFound
	}

	if post.UserID != userID {
		return ErrForbidden
	}

	if err = s.statService.DecrementPosts(ctx, userID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, postID)
}
