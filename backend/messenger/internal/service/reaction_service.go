package service

import (
	"context"
	"encoding/json"
	"messenger/internal/errors"
	"messenger/internal/model"
	"shared/events"
	"time"

	"github.com/google/uuid"
)

func (s *ChatService) AddReaction(ctx context.Context, messageID, userID, requesterID uuid.UUID, emoji string) error {
	if userID != requesterID {
		return errors.ErrChatAccessDenied
	}

	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.ErrMessageNotFound
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
		return err
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
		return errors.ErrChatAccessDenied
	}

	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.ErrMessageNotFound
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
		return err
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeReactionRemoved,
		Timestamp: time.Now(),
		Data:      payload,
	}

	return s.publishEvent(event)
}
