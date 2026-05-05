package postgres

import (
	"context"
	"errors"

	"social_network/internal/comment"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	*BaseRepository
}

func NewCommentRepository(db *gorm.DB) (comment.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &CommentRepository{BaseRepository: repo}, nil
}

func (r *CommentRepository) Create(ctx context.Context, comment *comment.Comment) error {
	return r.DB(ctx).Create(comment).Error
}

func (r *CommentRepository) GetByID(ctx context.Context, commentID uuid.UUID) (*comment.Comment, error) {
	var comment comment.Comment

	err := r.DB(ctx).First(&comment, "id = ?", commentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}

		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) ListByPostID(ctx context.Context, postID uuid.UUID) ([]comment.Comment, error) {
	var comments []comment.Comment
	err := r.DB(ctx).Where("post_id = ?", postID).Order("created_at asc").Find(&comments).Error

	return comments, err
}

func (r *CommentRepository) Delete(ctx context.Context, commentID, userID uuid.UUID) error {
	deleted := r.DB(ctx).Where("id = ? AND user_id = ?", commentID, userID).Delete(&comment.Comment{})

	if deleted.RowsAffected == 0 {
		return comment.ErrCommentNotFound
	}

	return deleted.Error
}
