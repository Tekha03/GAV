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

    payload, _ := json.Marshal(events.ReactionAddedData{
        MessageID: messageID,
        UserID: userID,
        Reaction: emoji,
    })

    event := events.Event{
        EventID: uuid.New(),
        EventType: events.EventTypeReactionAdded,
        Timestamp: time.Now(),
        Data: payload,
    }

    if err := s.producer.PublishEvent(event); err != nil {
        return err
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

    payload, _ := json.Marshal(events.ReactionRemovedData{
        MessageID: messageID,
        UserID: userID,
    })


    event := events.Event{
        EventID: uuid.New(),
        EventType: events.EventTypeReactionRemoved,
        Timestamp: time.Now(),
        Data: payload,
    }

    if err := s.producer.PublishEvent(event); err != nil {
        return err
    }

    err = s.reactionRepo.Remove(ctx, messageID, userID)
    return err
}
