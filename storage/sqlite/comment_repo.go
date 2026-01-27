package sqlite

import (
	"gav/internal/comment"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (cr *CommentRepository) Create(comment *comment.Comment) error {
	return cr.db.Create(comment).Error
}
