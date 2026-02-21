package post

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, postID uuid.UUID) (*Post, error)
	ListByUser(ctx context.Context, authorID uuid.UUID) ([]*Post, error)
	Delete(ctx context.Context, postID uuid.UUID) error
}
