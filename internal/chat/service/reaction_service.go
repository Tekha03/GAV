package service

import (
	"context"
	"gav/internal/chat/errors"
	"gav/internal/chat/model"

	"github.com/google/uuid"
)

func (s *ChatService) AddReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
    msg, err := s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return err
    }

    if msg == nil {
        return errors.ErrMessageNotFound
    }

    reaction := &model.Reaction{
        ID: uuid.New(),
        MessageID: messageID,
        UserID: userID,
        Emoji: emoji,
    }

    err = s.reactionRepo.Add(ctx, reaction)

    return err
}
func (s *ChatService) RemoveReaction(ctx context.Context, messageID, userID uuid.UUID) error {
    msg, err := s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return err
    }

    if msg == nil {
        return errors.ErrMessageNotFound
    }

    err = s.reactionRepo.Remove(ctx, messageID, userID)

    return err
}