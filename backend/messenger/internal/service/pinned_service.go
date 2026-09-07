package service

import (
	"context"
	"messenger/internal/model"
	apperrors "shared/app_errors"

	"github.com/google/uuid"
)

func (s *ChatService) PinMessage(ctx context.Context, requesterID, messageID uuid.UUID) error {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return apperrors.New(apperrors.MessageNotFound, "message not found")
	}
	if err := s.RequireChatMember(ctx, msg.ChatID, requesterID); err != nil {
		return err
	}

	return s.pinnedRepo.Pin(ctx, msg.ChatID, messageID)
}

func (s *ChatService) UnpinMessage(ctx context.Context, requesterID, messageID uuid.UUID) error {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return apperrors.New(apperrors.MessageNotFound, "message not found")
	}
	if err := s.RequireChatMember(ctx, msg.ChatID, requesterID); err != nil {
		return err
	}

	return s.pinnedRepo.Unpin(ctx, msg.ChatID, messageID)
}

func (s *ChatService) GetPinnedMessages(ctx context.Context, chatID, requesterID uuid.UUID) ([]*model.Message, error) {
	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return nil, err
	}

	ids, err := s.pinnedRepo.GetByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
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
