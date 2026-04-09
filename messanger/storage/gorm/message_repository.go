// messanger/chat/storage/gorm/message_repository.go
package gorm

import (
	"context"
	"errors"
	"messanger/chat/model"
	"messanger/chat/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
    repo *Repository
}

func NewMessageRepository(repo *Repository) repository.MessageRepository {
    return &MessageRepository{repo: repo}
}

func (mr *MessageRepository) Create(ctx context.Context, msg *model.Message) (uuid.UUID, error) {
    result := mr.repo.WithContext(ctx).Create(msg)
    if result.Error != nil {
        return uuid.Nil, result.Error
    }
    return msg.ID, nil
}

func (mr *MessageRepository) UpdateText(ctx context.Context, messageID uuid.UUID, newText string) error {
    return mr.repo.WithContext(ctx).
        Model(&model.Message{}).
        Where("id = ? AND deleted_at IS NULL", messageID).
        Updates(map[string]interface{}{
            "text":      newText,
            "edited_at": time.Now(),
        }).Error
}

func (mr *MessageRepository) Delete(ctx context.Context, messageID uuid.UUID) error {
    return mr.repo.WithContext(ctx).
        Model(&model.Message{}).
        Where("id = ?", messageID).
        Updates(map[string]interface{}{
            "deleted_at": time.Now(),
        }).Error
}

func (mr *MessageRepository) GetByID(ctx context.Context, messageID uuid.UUID) (*model.Message, error) {
    var msg model.Message
    err := mr.repo.WithContext(ctx).
        Where("id = ? AND deleted_at IS NULL", messageID).
        First(&msg).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrMessageNotFound
        }
        return nil, err
    }
    return &msg, nil
}

func (mr *MessageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
    query := mr.repo.WithContext(ctx).
        Where("chat_id = ? AND deleted_at IS NULL", chatID).
        Order("created_at DESC").
        Limit(limit)
    
    if cursorID != nil {
        query = query.Where("created_at < ?", *cursorID)
    }
    
    var messages []*model.Message
    return messages, query.Find(&messages).Error
}

func (mr *MessageRepository) UpdateReadAtForChat(ctx context.Context, chatID, userID uuid.UUID, readAt time.Time) error {
    return nil
}
