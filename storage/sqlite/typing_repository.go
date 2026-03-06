package sqlite

import (
	"context"
	"errors"
	"gav/internal/chat/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *MessageRepository) Forward(ctx context.Context, messageID, targetChatID, senderID uuid.UUID) (*model.Message, error) {
	var orig model.Message

	if err := r.DB(ctx).First(&orig, "id = ?", messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("original message not found")
		}
		return nil, err
	}

	newMsg := model.Message{
		ID:        uuid.New(),
		ChatID:    targetChatID,
		SenderID:  senderID,
		Text:      orig.Text,
		ReplyToID: nil,
		CreatedAt: time.Now(),
	}

	if err := r.DB(ctx).Create(&newMsg).Error; err != nil {
		return nil, err
	}

	var attachments []model.Attachment
	if err := r.DB(ctx).Where("message_id = ?", orig.ID).Find(&attachments).Error; err != nil {
		return nil, err
	}

	for _, att := range attachments {
		newAtt := model.Attachment{
			ID:        uuid.New(),
			MessageID: newMsg.ID,
			URL:       att.URL,
			Type:      att.Type,
			FileName:  att.FileName,
			FileSize:  att.FileSize,
		}
		if err := r.DB(ctx).Create(&newAtt).Error; err != nil {
			return nil, err
		}
	}

	return &newMsg, nil
}