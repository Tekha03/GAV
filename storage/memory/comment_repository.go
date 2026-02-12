package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"gav/internal/comment"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentRepository struct {
	mu       sync.RWMutex
	comments map[uint]comment.Comment
	nextID   uint
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		comments: make(map[uint]comment.Comment),
		nextID:   1,
	}
}

func (cr *CommentRepository) Create(ctx context.Context, comment *comment.Comment) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	comment.ID = cr.nextID
	comment.CreatedAt = time.Now()
	cr.comments[cr.nextID] = *comment
	cr.nextID++

	return nil
}

func (cr *CommentRepository) GetByPostID(ctx context.Context, postID uint) ([]comment.Comment, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var result []comment.Comment
	for _, comment := range cr.comments {
		if comment.PostID == postID {
			result = append(result, comment)
		}
	}

	return result, nil
}

func (cr *CommentRepository) Delete(ctx context.Context, commentID, userID uint) error {
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
