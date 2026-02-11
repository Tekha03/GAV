package post

import "context"

type Repository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, postID uint) (*Post, error)
	ListByUser(ctx context.Context, authorID uint) ([]*Post, error)
	Delete(ctx context.Context, postID uint) error
}
