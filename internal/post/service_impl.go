package post

import (
	"context"
	"errors"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden	= errors.New("forbidden")
	ErrEmptyContent = errors.New("empty content")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	userID uint,
	content string,
) (*Post, error) {

	if content == "" {
		return  nil, ErrEmptyContent
	}

	post := &Post{
		UserID: 	userID,
		Content: 	content,
		CreatedAt: 	time.Now(),
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) GetByID(ctx context.Context, postID uint) (*Post, error) {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (s *Service) ListByUser(ctx context.Context, userID uint) ([]*Post, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) Delete(ctx context.Context, userID, postID uint) error {
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

	return s.repo.Delete(ctx, postID)
}
