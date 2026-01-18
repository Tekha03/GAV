package sqlite

import (
	"context"
	"errors"

	"gav/post"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (pr *PostRepository) Create(ctx context.Context, post *post.Post) error {
	return pr.db.WithContext(ctx).Create(post).Error
}

func (pr *PostRepository) GetByID(ctx context.Context, id uint) (*post.Post, error) {
	var post *post.Post

	if err := pr.db.WithContext(ctx).First(&post, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}

		return nil, err
	}

	return post, nil
}

func (pr *PostRepository) ListByUser(ctx context.Context, userID uint) ([]*post.Post, error) {
	var posts []*post.Post

	if err := pr.db.WithContext(ctx).Where("user id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (pr *PostRepository) Delete(ctx context.Context, id uint) error {
	deleted := pr.db.WithContext(ctx).Delete(&post.Post{}, "id = ?", id)

	if deleted.Error != nil {
		return deleted.Error
	}

	if deleted.RowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}
