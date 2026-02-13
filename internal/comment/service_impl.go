package comment

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrForbiddenDelete = errors.New("forbidden: cannot delete someone else's comment")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID, postID uuid.UUID, content string) error {
	comment := &Comment{
		UserID: userID,
		PostID: postID,
		Content: content,
	}

	return  s.repo.Create(ctx, comment)
}

func (s *Service) GetByID(ctx context.Context, commentID uuid.UUID) (*Comment, error) {
	return s.repo.GetByID(ctx, commentID)
}

func (s *Service) GetByPostID(ctx context.Context, postID uuid.UUID) ([]Comment, error) {
	return s.repo.GetByPostID(ctx, postID)
}

func (s *Service) Delete(ctx context.Context, userID, commentID uuid.UUID) error {
	return s.repo.Delete(ctx, userID, commentID)
}
