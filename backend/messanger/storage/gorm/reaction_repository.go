package gorm

import (
	"context"
	"messanger/internal/model"
	"messanger/internal/repository"

	"github.com/google/uuid"
)

type ReactionRepository struct {
    repo *Repository
}

func NewReactionRepository(repo *Repository) repository.ReactionRepository {
    return &ReactionRepository{repo: repo}
}

func (rr *ReactionRepository) Add(ctx context.Context, reaction *model.Reaction) error {
    return rr.repo.WithContext(ctx).Create(reaction).Error
}

func (rr *ReactionRepository) Remove(ctx context.Context, messageID, userID uuid.UUID) error {
    return rr.repo.WithContext(ctx).
        Where("message_id = ? AND user_id = ?", messageID, userID).
        Delete(&model.Reaction{}).Error
}
