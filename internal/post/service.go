package post

import (
	"context"

	"github.com/google/uuid"
)

type PostService interface {
	Create(ctx context.Context, userID uuid.UUID, content string) (*Post, error)
	GetByID(ctx context.Context, postID uuid.UUID) (*Post, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	Delete(ctx context.Context, userID, postID uuid.UUID) error
}
