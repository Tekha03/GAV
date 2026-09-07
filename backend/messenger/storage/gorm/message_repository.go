// messanger/chat/storage/gorm/message_repository.go
package gorm

import (
	"context"
	"errors"
	"messenger/internal/model"
	"messenger/internal/repository"
	apperrors "shared/app_errors"
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
		return uuid.Nil, createError(result.Error, apperrors.MessageAlreadyExists, "message already exists", "failed to create message")
	}
	return msg.ID, nil
}

func (mr *MessageRepository) UpdateText(ctx context.Context, messageID uuid.UUID, newText string) error {
	result := mr.repo.WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ? AND deleted_at IS NULL", messageID).
		Updates(map[string]interface{}{
			"text":      newText,
			"edited_at": time.Now(),
		})
	return mutationError(result, apperrors.MessageNotFound, "message not found", "failed to update message text")
}

func (mr *MessageRepository) Delete(ctx context.Context, messageID uuid.UUID) error {
	result := mr.repo.WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
		})
	return mutationError(result, apperrors.MessageNotFound, "message not found", "failed to delete message")
}

func (mr *MessageRepository) GetByID(ctx context.Context, messageID uuid.UUID) (*model.Message, error) {
	var msg model.Message
	err := mr.repo.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", messageID).
		First(&msg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.MessageNotFound, "message not found")
		}
		return nil, internalError("failed to get message by ID", err)
	}
	return &msg, nil
}

func (mr *MessageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
	query := mr.repo.WithContext(ctx).
		Where("chat_id = ? AND deleted_at IS NULL", chatID).
		Order("created_at DESC, id DESC").
		Limit(limit)

	if cursorID != nil {
		var cursorMessage model.Message
		err := mr.repo.WithContext(ctx).
			Select("id", "created_at").
			Where("id = ? AND chat_id = ? AND deleted_at IS NULL", *cursorID, chatID).
			First(&cursorMessage).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperrors.New(apperrors.MessageNotFound, "cursor message not found")
			}
			return nil, internalError("failed to get cursor message", err)
		}

		query = query.Where(
			"(created_at < ?) OR (created_at = ? AND id < ?)",
			cursorMessage.CreatedAt,
			cursorMessage.CreatedAt,
			cursorMessage.ID,
		)
	}

	var messages []*model.Message
	if err := query.Find(&messages).Error; err != nil {
		return nil, internalError("failed to get messages by chat ID", err)
	}
	return messages, nil
}

func (mr *MessageRepository) UpdateLastReadMessageForChat(ctx context.Context, chatID, userID uuid.UUID) error {
	var lastMessage model.Message
	lastMessageID := uuid.Nil

	err := mr.repo.WithContext(ctx).
		Select("id").
		Where("chat_id = ? AND deleted_at IS NULL", chatID).
		Order("created_at DESC, id DESC").
		Limit(1).
		First(&lastMessage).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return internalError("failed to get last chat message", err)
		}
	} else {
		lastMessageID = lastMessage.ID
	}

	result := mr.repo.WithContext(ctx).
		Model(&model.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("last_read_message_id", lastMessageID)

	return mutationError(result, apperrors.ChatMemberNotFound, "chat member not found", "failed to update last read message")
}
