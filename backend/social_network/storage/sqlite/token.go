package sqlite

import (
	"context"
	"errors"
	"social_network/internal/token"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository struct {
	*BaseRepository
}

func NewTokenRepository(db *gorm.DB) (token.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &TokenRepository{BaseRepository: repo}, nil
}

func (r *TokenRepository) Create(ctx context.Context, token *token.RefreshToken) error {
	return r.DB(ctx).Create(token).Error
}

func (r *TokenRepository) GetByHash(ctx context.Context, hash string) (*token.RefreshToken, error) {
	var refreshToken token.RefreshToken

	if err := r.DB(ctx).Where("token_hash = ?", hash).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenNotFound
		}

		return nil, err
	}

	return &refreshToken, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, hash string) error {
	result := r.DB(ctx).
		Model(&token.RefreshToken{}).
		Where("token_hash = ?", hash).
		Updates(map[string]interface{}{
			"revoked":    true,
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

func (r *TokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	result := r.DB(ctx).
		Model(&token.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":    true,
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
