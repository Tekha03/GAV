package comment

import (
	"context"
	"errors"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrForbiddenDelete = errors.New("forbidden: cannot delete someone else's comment")
)

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (cs *CommentService) Create(ctx context.Context, userID, postID uint, content string) error {
	comment := &Comment{
		UserID: userID,
		PostID: postID,
		Content: content,
	}

	return  cs.repo.Create(ctx, comment)
}

func (cs *CommentService) GetByPostID(ctx context.Context, postID uint) ([]Comment, error) {
	return cs.repo.GetByPostID(ctx, postID)
}

func (cs *CommentService) Delete(ctx context.Context, userID, commentID uint) error {
	return cs.repo.Delete(ctx, userID, commentID)
}
