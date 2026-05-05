package postgres

import (
	"context"
	"errors"
	"time"

	"social_network/internal/post"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	*BaseRepository
}

func NewPostRepository(db *gorm.DB) (post.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &PostRepository{BaseRepository: repo}, nil
}

func (r *PostRepository) Create(ctx context.Context, post *post.Post) error {
	return r.DB(ctx).Create(post).Error
}

func (r *PostRepository) GetByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	var post *post.Post

	if err := r.DB(ctx).First(&post, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}

		return nil, err
	}

	return post, nil
}

func (r *PostRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*post.Post, error) {
	var posts []*post.Post

	if err := r.DB(ctx).Where("user id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) ListFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*post.Post, error) {
	var posts []*post.Post

	query := r.DB(ctx).
		Table("posts").
		Joins("JOIN follows ON posts.user_id = follows.following_id").
		Where("follows.following_id = ?", userID)

	if !before.IsZero() {
		query = query.Where("posts.created_at < ?", before)
	}

	err := query.Order("posts.created_at DESC, posts.id DESC").Limit(limit + 1).Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	deleted := r.DB(ctx).Delete(&post.Post{}, "id = ?", id)

	if deleted.Error != nil {
		return deleted.Error
	}

	if deleted.RowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}
