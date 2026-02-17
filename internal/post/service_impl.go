package post

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden	= errors.New("forbidden")
	ErrEmptyContent = errors.New("empty content")
)

type service struct {
	repo Repository
}

func NewService(repo Repository) PostService {
	return &service{repo: repo}
}

func (s *service) Create(
	ctx context.Context,
	userID uuid.UUID,
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

	return s.repo.Delete(ctx, postID)
}
