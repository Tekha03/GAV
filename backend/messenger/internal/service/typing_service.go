package service

import (
	"context"
	"messenger/internal/errors"

	"github.com/google/uuid"
)

func (s *ChatService) SendTyping(ctx context.Context, chatID, userID uuid.UUID) error {
	members, err := s.membersRepo.GetMembers(ctx, chatID)
	if err != nil {
		return err
	}

	found := false
	for _, m := range members {
		if m.UserID == userID {
			found = true
			break
		}
	}
	if !found {
		return errors.ErrMemberNotFound
	}

	s.typingRepo.SetTyping(ctx, chatID, userID)
	return nil
}
