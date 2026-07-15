package like

import (
	"context"
	"social_network/internal/stats"

	"github.com/google/uuid"
)

type service struct {
	repo        Repository
	statService stats.StatsService
}

func NewService(repo Repository, statService ...stats.StatsService) (LikeService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo, statService: stats.ServiceOrNoop(statService...)}, nil
}

func (s *service) Add(ctx context.Context, like Like) error {
	if like.UserID == uuid.Nil || like.PostID == uuid.Nil {
		return ErrInvalidLike
	}

	alreadyLiked, err := s.repo.LikeExists(ctx, like)
	if err != nil {
		return err
	}

	if alreadyLiked {
		return ErrAlreadyLiked
	}

	if err = s.statService.IncrementPostLikes(ctx, like.PostID); err != nil {
		return err
	}

	return s.repo.Add(ctx, like)
}

func (s *service) Remove(ctx context.Context, like Like) error {
	if like.UserID == uuid.Nil || like.PostID == uuid.Nil {
		return ErrInvalidLike
	}

	likeExists, err := s.repo.LikeExists(ctx, like)
	if err != nil {
		return err
	}

	if !likeExists {
		return ErrLikeDoesNotExist
	}

	if err = s.statService.DecrementPostLikes(ctx, like.PostID); err != nil {
		return err
	}

	return s.repo.Remove(ctx, like)
}
