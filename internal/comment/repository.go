package comment

import "context"

type Repository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByID(ctx context.Context, commentID uint) (*Comment, error)
	GetByPostID(ctx context.Context, postID uint) ([]Comment, error)
	Delete(ctx context.Context, userID, commentID uint) error
}
