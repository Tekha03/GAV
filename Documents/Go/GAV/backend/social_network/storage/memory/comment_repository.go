package memory

import (
	"context"
	"sync"
	"time"

	"social_network/internal/comment"

	"github.com/google/uuid"
)

type CommentRepository struct {
	mu			sync.RWMutex
	comments	map[uuid.UUID]comment.Comment
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		comments: make(map[uuid.UUID]comment.Comment),
	}
}

func (r *CommentRepository) Create(ctx context.Context, comment *comment.Comment) error {
	if comment == nil {
		return ErrCommentNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()
	r.comments[comment.ID] = *comment

	return nil
}

func (r *CommentRepository) ListyPostID(ctx context.Context, postID uuid.UUID) ([]comment.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []comment.Comment
	for _, comment := range r.comments {
		if comment.PostID == postID {
			result = append(result, comment)
		}
	}

	return result, nil
}

func (r *CommentRepository) Delete(ctx context.Context, commentID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	comment, ok := r.comments[commentID]
	if !ok {
		return ErrCommentNotFound
	}

	if comment.UserID != userID {
		return ErrCommentNotFound
	}

	delete(r.comments, commentID)
	return nil
}
