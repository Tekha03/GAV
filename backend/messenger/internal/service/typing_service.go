package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *ChatService) SendTyping(ctx context.Context, chatID, userID uuid.UUID) error {
	if err := s.requireChatMember(ctx, chatID, userID); err != nil {
		return err
	}

	s.typingRepo.SetTyping(ctx, chatID, userID)
	return nil
}
