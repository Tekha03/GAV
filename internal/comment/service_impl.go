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

type service struct {
	repo Repository
}

func NewService(repo Repository) CommentService {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID, postID uuid.UUID, content string) error {
	comment := &Comment{
		UserID: userID,
		PostID: postID,
		Content: content,
	}

	return  s.repo.Create(ctx, comment)
}

func (s *service) GetByID(ctx context.Context, commentID uuid.UUID) (*Comment, error) {
	return s.repo.GetByID(ctx, commentID)
}

func (s *service) GetByPostID(ctx context.Context, postID uuid.UUID) ([]Comment, error) {
	return s.repo.GetByPostID(ctx, postID)
}

func (s *service) Delete(ctx context.Context, userID, commentID uuid.UUID) error {
	return s.repo.Delete(ctx, userID, commentID)
}
