package service

import (
	"context"
	"encoding/json"
	"messenger/internal/model"
	apperrors "shared/app_errors"
	"shared/events"
	"time"

	"github.com/google/uuid"
)

func (s *ChatService) AddReaction(ctx context.Context, messageID, userID, requesterID uuid.UUID, emoji string) error {
	if userID != requesterID {
		return apperrors.New(apperrors.ChatAccessDenied, "chat access denied")
	}

	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return apperrors.New(apperrors.MessageNotFound, "message not found")
	}
	if err := s.requireChatMember(ctx, msg.ChatID, requesterID); err != nil {
		return err
	}

	reaction := &model.Reaction{
		ID:        uuid.New(),
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	}

	if err := s.reactionRepo.Add(ctx, reaction); err != nil {
		return err
	}

	payload, err := json.Marshal(events.ReactionAddedData{
		MessageID: messageID,
		UserID:    userID,
		Reaction:  emoji,
	})
	if err != nil {
		return apperrors.Wrap(apperrors.Internal, "failed to encode reaction added event", err)
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeReactionAdded,
		Timestamp: time.Now(),
		Data:      payload,
	}

	return s.publishEvent(event)
}

func (s *ChatService) RemoveReaction(ctx context.Context, messageID, userID, requesterID uuid.UUID) error {
	if userID != requesterID {
		return apperrors.New(apperrors.ChatAccessDenied, "chat access denied")
	}

	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return apperrors.New(apperrors.MessageNotFound, "message not found")
	}
	if err := s.requireChatMember(ctx, msg.ChatID, requesterID); err != nil {
		return err
	}

	if err := s.reactionRepo.Remove(ctx, messageID, userID); err != nil {
		return err
	}

	payload, err := json.Marshal(events.ReactionRemovedData{
		MessageID: messageID,
		UserID:    userID,
	})
	if err != nil {
		return apperrors.Wrap(apperrors.Internal, "failed to encode reaction removed event", err)
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeReactionRemoved,
		Timestamp: time.Now(),
		Data:      payload,
	}

	return s.publishEvent(event)
}
