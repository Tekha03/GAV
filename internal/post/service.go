package post

import (
	"context"
)

type PostService interface {
	Create(ctx context.Context, userID uint, content string) (*Post, error)
	GetByID(ctx context.Context, postID uint) (*Post, error)
	ListByUser(ctx context.Context, userID uint) ([]*Post, error)
	Delete(ctx context.Context, userID, postID uint) error
}
