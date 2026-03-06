package sqlite

import (
	"context"
	"errors"
	"gav/internal/chat/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReactionRepository struct {
	*BaseRepository
}

func NewReactionRepository(db *gorm.DB) (*ReactionRepository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &ReactionRepository{BaseRepository: repo}, nil
}

func (r *ReactionRepository) Add(ctx context.Context, reaction *model.Reaction) error {
	if reaction.ID == uuid.Nil {
		reaction.ID = uuid.New()
	}
	return r.DB(ctx).Create(reaction).Error
}

func (r *ReactionRepository) Remove(ctx context.Context, messageID, userID uuid.UUID) error {
	res := r.DB(ctx).
		Where("message_id = ? AND user_id = ?", messageID, userID).
		Delete(&model.Reaction{})

	if res.RowsAffected == 0 {
		return errors.New("reaction not found")
	}
	return res.Error
}

func (r *ReactionRepository) GetByMessage(ctx context.Context, messageID uuid.UUID) ([]*model.Reaction, error) {
	var reactions []*model.Reaction
	err := r.DB(ctx).
		Where("message_id = ?", messageID).
		Find(&reactions).Error
	if err != nil {
		return nil, err
	}
	return reactions, nil
}