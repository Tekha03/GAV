package comment

import "context"

type Service interface {
	Create(ctx context.Context, userID, postID uint, content string) error
	GetByPostID(ctx context.Context, postID uint) ([]Comment, error)
	Delete(ctx context.Context, userID, commentID uint) error
}
