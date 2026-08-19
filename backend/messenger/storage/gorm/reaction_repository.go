package gorm

import (
	"context"
	"messenger/internal/model"
	"messenger/internal/repository"
	apperrors "shared/app_errors"

	"github.com/google/uuid"
)

type ReactionRepository struct {
	repo *Repository
}

func NewReactionRepository(repo *Repository) repository.ReactionRepository {
	return &ReactionRepository{repo: repo}
}

func (rr *ReactionRepository) Add(ctx context.Context, reaction *model.Reaction) error {
	err := rr.repo.WithContext(ctx).Create(reaction).Error
	return createError(err, apperrors.ReactionAlreadyExists, "reaction already exists", "failed to add reaction")
}

func (rr *ReactionRepository) Remove(ctx context.Context, messageID, userID uuid.UUID) error {
	result := rr.repo.WithContext(ctx).
		Where("message_id = ? AND user_id = ?", messageID, userID).
		Delete(&model.Reaction{})
	return mutationError(result, apperrors.ReactionNotFound, "reaction not found", "failed to remove reaction")
}
