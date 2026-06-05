package service

import (
	"context"
	"messenger/internal/errors"
	"messenger/internal/model"

	"github.com/google/uuid"
)

func (s *ChatService) PinMessage(ctx context.Context, messageID uuid.UUID) error {
    msg, err := s.messageRepo.GetByID(ctx, messageID)
    if err != nil {
        return err
    }
    if msg == nil {
        return errors.ErrMessageNotFound
    }

    return s.pinnedRepo.Pin(ctx, msg.ChatID, messageID)
}

func (s *ChatService) UnpinMessage(ctx context.Context, messageID uuid.UUID) error {
    msg, err := s.messageRepo.GetByID(ctx, messageID)
    if err != nil {
        return err
    }
    if msg == nil {
        return errors.ErrMessageNotFound
    }

    return s.pinnedRepo.Unpin(ctx, msg.ChatID, messageID)
}

func (s *ChatService) GetPinnedMessages(ctx context.Context, chatID uuid.UUID) ([]*model.Message, error) {
    ids := s.pinnedRepo.GetByChatID(ctx, chatID)
    messages := []*model.Message{}

    for _, id := range ids {
        msg, err := s.messageRepo.GetByID(ctx, id)
        if err != nil {
            return nil, err
        }
        if msg != nil {
            messages = append(messages, msg)
        }
    }

    return messages, nil
}
