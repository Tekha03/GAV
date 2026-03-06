package sqlite

import (
	"context"
	"errors"
	"gav/internal/chat/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PinnedRepository struct {
	*BaseRepository
}

func NewPinnedRepository(db *gorm.DB) (*PinnedRepository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &PinnedRepository{BaseRepository: repo}, nil
}

func (r *PinnedRepository) Pin(ctx context.Context, chatID, messageID uuid.UUID) error {
	p := model.PinnedMessages{
		ChatID:    chatID,
		MessageID: messageID,
		PinnedAt:  time.Now(),
	}
	return r.DB(ctx).Create(&p).Error
}

func (r *PinnedRepository) Unpin(ctx context.Context, chatID, messageID uuid.UUID) error {
	res := r.DB(ctx).
		Where("chat_id = ? AND message_id = ?", chatID, messageID).
		Delete(&model.PinnedMessages{})

	if res.RowsAffected == 0 {
		return errors.New("pinned message not found")
	}
	return res.Error
}

func (r *PinnedRepository) GetByChatID(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	var pinned []model.PinnedMessages
	err := r.DB(ctx).
		Where("chat_id = ?", chatID).
		Find(&pinned).Error
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	for _, p := range pinned {
		ids = append(ids, p.MessageID)
	}

	return ids, nil
}