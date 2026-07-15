package sqlite

import (
	"context"
	"errors"
	"fmt"
	"social_network/internal/stats"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StatsRepository struct {
	*BaseRepository
}

func NewStatsRepository(db *gorm.DB) (stats.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &StatsRepository{BaseRepository: repo}, nil
}

func (r *StatsRepository) CreateUserStats(ctx context.Context, userStats *stats.UserStats) error {
	if err := r.DB(ctx).Create(userStats).Error; err != nil {
		return fmt.Errorf("stats repository: create user stats: %w", err)
	}

	return nil
}

func (r *StatsRepository) DeleteUserStats(ctx context.Context, userID uuid.UUID) error {
	result := r.DB(ctx).Delete(&stats.UserStats{}, "user_id = ?", userID)

	if result.Error != nil {
		return fmt.Errorf("stats repository: delete user stats: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return stats.ErrStatsNotFound
	}

	return nil
}

func (r *StatsRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*stats.UserStats, error) {
	var userStats stats.UserStats

	if err := r.DB(ctx).First(&userStats, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, stats.ErrStatsNotFound
		}

		return nil, fmt.Errorf("stats repository: get user stats: %w", err)
	}

	return &userStats, nil
}

func (r *StatsRepository) IncrementPosts(ctx context.Context, userID uuid.UUID) error {
	return r.incrementUserField(ctx, userID, "post_count", 1)
}

func (r *StatsRepository) DecrementPosts(ctx context.Context, userID uuid.UUID) error {
	return r.incrementUserField(ctx, userID, "post_count", -1)
}

func (r *StatsRepository) IncrementDogs(ctx context.Context, userID uuid.UUID) error {
	return r.incrementUserField(ctx, userID, "dogs_count", 1)
}

func (r *StatsRepository) DecrementDogs(ctx context.Context, userID uuid.UUID) error {
	return r.incrementUserField(ctx, userID, "dogs_count", -1)
}

func (r *StatsRepository) IncrementFollowers(ctx context.Context, userID uuid.UUID) error {
	return r.incrementUserField(ctx, userID, "followers", 1)
}

func (r *StatsRepository) DecrementFollowers(ctx context.Context, userID uuid.UUID) error {
	return r.incrementUserField(ctx, userID, "followers", -1)
}

func (r *StatsRepository) IncrementFollowings(ctx context.Context, userID uuid.UUID) error {
	return r.incrementUserField(ctx, userID, "followings", 1)
}

func (r *StatsRepository) DecrementFollowings(ctx context.Context, userID uuid.UUID) error {
	return r.incrementUserField(ctx, userID, "followings", -1)
}

func (r *StatsRepository) incrementUserField(ctx context.Context, userID uuid.UUID, field string, delta int) error {
	var userStats stats.UserStats
	if err := r.DB(ctx).
		Where("user_id = ?", userID).
		FirstOrCreate(&userStats, stats.UserStats{UserID: userID}).Error; err != nil {
		return fmt.Errorf("stats repository: ensure user stats: %w", err)
	}

	result := r.DB(ctx).
		Model(&stats.UserStats{}).
		Where("user_id = ?", userID).UpdateColumn(field, gorm.Expr(field+" + ?", delta))

	if result.Error != nil {
		return fmt.Errorf("stats repository: increment user field: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return stats.ErrStatsNotFound
	}

	return nil
}

func (r *StatsRepository) CreatePostStats(ctx context.Context, postStats *stats.PostStats) error {
	if err := r.DB(ctx).Create(postStats).Error; err != nil {
		return fmt.Errorf("stats repository: create post stats: %w", err)
	}

	return nil
}

func (r *StatsRepository) GetPostStats(ctx context.Context, postID uuid.UUID) (*stats.PostStats, error) {
	var postStats stats.PostStats

	if err := r.DB(ctx).First(&postStats, "post_id = ?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, stats.ErrStatsNotFound
		}

		return nil, fmt.Errorf("stats repository: get post stats: %w", err)
	}

	return &postStats, nil
}

func (r *StatsRepository) IncrementPostLikes(ctx context.Context, postID uuid.UUID) error {
	return r.incrementPostField(ctx, postID, "likes_count", 1)
}

func (r *StatsRepository) DecrementPostLikes(ctx context.Context, postID uuid.UUID) error {
	return r.incrementPostField(ctx, postID, "likes_count", -1)
}

func (r *StatsRepository) IncrementPostComments(ctx context.Context, postID uuid.UUID) error {
	return r.incrementPostField(ctx, postID, "comments_count", 1)
}

func (r *StatsRepository) DecrementPostComments(ctx context.Context, postID uuid.UUID) error {
	return r.incrementPostField(ctx, postID, "comments_count", -1)
}

func (r *StatsRepository) incrementPostField(ctx context.Context, postID uuid.UUID, field string, delta int) error {
	var postStats stats.PostStats
	if err := r.DB(ctx).
		Where("post_id = ?", postID).
		FirstOrCreate(&postStats, stats.PostStats{PostID: postID}).Error; err != nil {
		return fmt.Errorf("stats repository: ensure post stats: %w", err)
	}

	result := r.DB(ctx).
		Model(&stats.PostStats{}).
		Where("post_id = ?", postID).
		UpdateColumn(field, gorm.Expr(field+" + ?", delta))

	if result.Error != nil {
		return fmt.Errorf("stats repository: increment post field: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return stats.ErrStatsNotFound
	}

	return nil
}
