package like

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) (LikeService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo}, nil
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

	return s.repo.Remove(ctx, like)
}
