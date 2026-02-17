package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"gav/internal/comment"

	"github.com/google/uuid"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentRepository struct {
	mu			sync.RWMutex
	comments	map[uuid.UUID]comment.Comment
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		comments: make(map[uuid.UUID]comment.Comment),
	}
}

func (cr *CommentRepository) Create(ctx context.Context, comment *comment.Comment) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()
	cr.comments[comment.ID] = *comment

	return nil
}

func (cr *CommentRepository) GetByPostID(ctx context.Context, postID uuid.UUID) ([]comment.Comment, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	var result []comment.Comment
	for _, comment := range cr.comments {
		if comment.PostID == postID {
			result = append(result, comment)
		}
	}

	return result, nil
}

func (cr *CommentRepository) Delete(ctx context.Context, commentID, userID uuid.UUID) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	comment, ok := cr.comments[commentID]
	if !ok {
		return ErrCommentNotFound
	}

	if comment.UserID != userID {
		return ErrCommentNotFound
	}

	delete(cr.comments, commentID)
	return nil
}
