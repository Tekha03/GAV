package comment

import "context"

type CommentService interface{
	Create(ctx context.Context, userID, postID uint, content string) error
	GetByID(ctx context.Context, commentID uint) (*Comment, error)
	GetByPostID(ctx context.Context, postID uint) ([]Comment, error)
	Delete(ctx context.Context, userID, commentID uint) error
}
