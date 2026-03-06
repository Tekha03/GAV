package sqlite

import (
	"context"
	"time"

	"gav/internal/chat/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	*BaseRepository
}

func NewMessageRepository(db *gorm.DB) (*MessageRepository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &MessageRepository{BaseRepository: repo}, nil
}

func (mr *MessageRepository) Create(ctx context.Context, msg *model.Message) (uuid.UUID, error) {
	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}

	msg.CreatedAt = time.Now()

	if err := mr.DB(ctx).Create(msg).Error; err != nil {
		return uuid.Nil, err
	}

	return msg.ID, nil
}

func (r *MessageRepository) UpdateText(ctx context.Context, chatID, msgID uuid.UUID, newText string) error {
	now := time.Now()

	return r.DB(ctx).
		Model(&model.Message{}).
		Where("id = ? AND chat_id = ?", msgID, chatID).
		Updates(map[string]interface{}{
			"text":      newText,
			"edited_at": &now,
		}).Error
}

func (r *MessageRepository) Delete(ctx context.Context, chatID, msgID uuid.UUID) error {
	now := time.Now()

	return r.DB(ctx).
		Model(&model.Message{}).
		Where("id = ? AND chat_id = ?", msgID, chatID).
		Update("deleted_at", &now).Error
}

func (r *MessageRepository) GetByID(ctx context.Context, chatID, msgID uuid.UUID) (*model.Message, error) {
	var msg model.Message

	err := r.DB(ctx).
		Where("id = ? AND chat_id = ?", msgID, chatID).
		First(&msg).Error

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (r *MessageRepository) GetByChatID(
	ctx context.Context,
	chatID uuid.UUID,
	limit int,
	cursorID *uuid.UUID,
) ([]*model.Message, error) {

	db := r.DB(ctx).
		Where("chat_id = ?", chatID).
		Where("deleted_at IS NULL")

	if cursorID != nil {
		var cursor model.Message

		if err := r.DB(ctx).
			Select("created_at").
			Where("id = ?", *cursorID).
			First(&cursor).Error; err != nil {
			return nil, err
		}

		db = db.Where("created_at < ?", cursor.CreatedAt)
	}

	var msgs []*model.Message

	err := db.
		Order("created_at DESC").
		Limit(limit).
		Find(&msgs).Error

	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *MessageRepository) UpdateReadAtForChat(
	ctx context.Context,
	chatID, userID uuid.UUID,
	readAt time.Time,
) error {

	return r.DB(ctx).
		Model(&model.Message{}).
		Where("chat_id = ?", chatID).
		Where("sender_id != ?", userID).
		Update("read_at", readAt).Error
}