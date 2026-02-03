package sqlite

import (
	"context"

	"gav/internal/comment"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (cr *CommentRepository) Create(ctx context.Context, comment *comment.Comment) error {
	return cr.db.WithContext(ctx).Create(comment).Error
}

func (cr *CommentRepository) GetByPostID(ctx context.Context, postID uint) ([]comment.Comment, error) {
	var comments []comment.Comment
	err := cr.db.WithContext(ctx).Where("post_id = ?", postID).Order("created_at asc").Find(&comments).Error

	return comments, err
}

func (cr *CommentRepository) Delete(ctx context.Context, commentID, userID uint) error {
	deleted := cr.db.WithContext(ctx).Where("id = ? AND user_id = ?", commentID, userID).Delete(&comment.Comment{})

	if deleted.RowsAffected == 0 {
		return comment.ErrCommentNotFound
	}

	return deleted.Error
}
