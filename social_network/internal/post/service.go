package post

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PostService interface {
	Create(ctx context.Context, userID uuid.UUID, content, imageUrl string) (*Post, error)
	GetByID(ctx context.Context, postID uuid.UUID) (*Post, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*Post, time.Time, error)
	Delete(ctx context.Context, userID, postID uuid.UUID) error
}
