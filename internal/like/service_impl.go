package like

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidLike = errors.New("invalid like")
	ErrAlreadyLiked = errors.New("already liked")
	ErrLikeDoesNotExist = errors.New("like does not exist")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(ctx context.Context, like Like) error {
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

func (s *Service) Remove(ctx context.Context, like Like) error {
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
