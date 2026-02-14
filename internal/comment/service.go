package comment

import (
	"context"

	"github.com/google/uuid"
)

type CommentService interface{
	Create(ctx context.Context, userID, postID uuid.UUID, content string) error
	GetByID(ctx context.Context, commentID uuid.UUID) (*Comment, error)
	GetByPostID(ctx context.Context, postID uuid.UUID) ([]Comment, error)
	Delete(ctx context.Context, userID, commentID uuid.UUID) error
}
