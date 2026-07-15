package service

import (
	"context"
	"encoding/json"
	"messenger/internal/constatnts"
	"messenger/internal/errors"
	"messenger/internal/model"
	"shared/events"
	"time"

	"github.com/google/uuid"
)

func (s *ChatService) SendMessage(ctx context.Context, requesterID uuid.UUID, input model.SendMessageInput) (*model.Message, error) {
	chat, err := s.chatRepo.GetByID(ctx, input.ChatID)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, errors.ErrChatNotFound
	}

	if input.SenderID != requesterID {
		return nil, errors.ErrChatAccessDenied
	}
	if err := s.requireChatMember(ctx, input.ChatID, requesterID); err != nil {
		return nil, err
	}

	if input.Text == nil && len(input.Attachments) == 0 {
		return nil, errors.ErrEmptyMessage
	}
	if input.Text != nil && len(*input.Text) > constatnts.MaxMessageLength {
		return nil, errors.ErrTextOverLength
	}
	if len(input.Attachments) > constatnts.MaxAttachments {
		return nil, errors.ErrAttachmentsOverLength
	}
	if input.ReplyToID != nil {
		msg, err := s.messageRepo.GetByID(ctx, *input.ReplyToID)
		if err != nil {
			return nil, err
		}
		if msg == nil || msg.ChatID != input.ChatID {
			return nil, errors.ErrInvalidReply
		}
	}

	message := &model.Message{
		ID:        uuid.New(),
		ChatID:    input.ChatID,
		SenderID:  input.SenderID,
		Text:      input.Text,
		ReplyToID: input.ReplyToID,
		CreatedAt: time.Now(),
	}

	msgID, err := s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, err
	}
	message.ID = msgID

	for _, att := range input.Attachments {
		attachment := model.Attachment{
			MessageID: msgID,
			URL:       att.URL,
			Type:      att.Type,
			FileName:  att.FileName,
			FileSize:  att.FileSize,
		}

		if err := s.attachmentRepo.Create(ctx, &attachment); err != nil {
			return nil, err
		}
		message.Attachments = append(message.Attachments, attachment)
	}

	receiverID, err := s.findChatReceiver(ctx, input.ChatID, input.SenderID)
	if err == nil && s.notClient != nil {
		var senderName string
		if s.socialClient != nil {
			usr, err := s.socialClient.GetUserProfile(ctx, input.SenderID)
			if err == nil && usr != nil {
				senderName = usr.Username
			}

			text := messageText(input.Text)

			go func() {
				_ = s.notClient.SendNewMessage(
					context.Background(),
					receiverID,
					senderName,
					text,
					input.ChatID.String(),
				)
			}()
		}
	}

	payload, err := json.Marshal(events.MessageSentData{
		MessageID: msgID,
		ChatID:    input.ChatID,
		SenderID:  input.SenderID,
		Text:      messageText(input.Text),
	})
	if err != nil {
		return nil, err
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeMessageSent,
		Timestamp: time.Now(),
		Data:      payload,
	}

	if err := s.publishEvent(event); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ChatService) EditMessage(ctx context.Context, requesterID, messageID uuid.UUID, newText string) (*model.Message, error) {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, errors.ErrMessageNotFound
	}
	if message.SenderID != requesterID {
		return nil, errors.ErrChatAccessDenied
	}
	if err := s.requireChatMember(ctx, message.ChatID, requesterID); err != nil {
		return nil, err
	}

	if err := s.messageRepo.UpdateText(ctx, messageID, newText); err != nil {
		return nil, err
	}

	message, err = s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, errors.ErrMessageNotFound
	}

	payload, err := json.Marshal(events.MessageEditedData{
		MessageID: message.ID,
		ChatID:    message.ChatID,
		Text:      messageText(message.Text),
	})
	if err != nil {
		return nil, err
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeMessageEdited,
		Timestamp: time.Now(),
		Data:      payload,
	}

	if err := s.publishEvent(event); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ChatService) DeleteMessage(ctx context.Context, requesterID, messageID uuid.UUID) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.ErrMessageNotFound
	}
	if message.SenderID != requesterID {
		return errors.ErrChatAccessDenied
	}
	if err := s.requireChatMember(ctx, message.ChatID, requesterID); err != nil {
		return err
	}

	if err := s.messageRepo.Delete(ctx, messageID); err != nil {
		return err
	}

	payload, err := json.Marshal(events.MessageDeletedData{
		MessageID: message.ID,
		ChatID:    message.ChatID,
	})
	if err != nil {
		return err
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeMessageDeleted,
		Timestamp: time.Now(),
		Data:      payload,
	}

	return s.publishEvent(event)
}

func (s *ChatService) GetMessages(ctx context.Context, chatID, requesterID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, errors.ErrChatNotFound
	}
	if err := s.requireChatMember(ctx, chatID, requesterID); err != nil {
		return nil, err
	}

	messages, err := s.messageRepo.GetByChatID(ctx, chatID, limit, cursorID)
	if err != nil {
		return nil, err
	}
	if messages == nil {
		return nil, errors.ErrMessageNotFound
	}

	for _, message := range messages {
		attachments, err := s.attachmentRepo.GetByMessage(ctx, message.ID)
		if err != nil {
			return nil, err
		}
		for _, attachment := range attachments {
			message.Attachments = append(message.Attachments, *attachment)
		}
	}

	return messages, nil
}

func (s *ChatService) MarkAsRead(ctx context.Context, chatID, requesterID uuid.UUID) error {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.ErrChatNotFound
	}
	if err := s.requireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	return s.messageRepo.UpdateReadAtForChat(ctx, chatID, requesterID, time.Now())
}

func (s *ChatService) ForwardMessage(ctx context.Context, requesterID, messageID, targetChatID uuid.UUID) (*model.Message, error) {
	origMsg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if origMsg == nil {
		return nil, errors.ErrMessageNotFound
	}

	if err := s.requireChatMember(ctx, origMsg.ChatID, requesterID); err != nil {
		return nil, err
	}
	if err := s.requireChatMember(ctx, targetChatID, requesterID); err != nil {
		return nil, err
	}

	input := model.SendMessageInput{
		ChatID:      targetChatID,
		SenderID:    requesterID,
		Text:        origMsg.Text,
		Attachments: []model.AttachmentInput{},
	}

	attachments, err := s.attachmentRepo.GetByMessage(ctx, origMsg.ID)
	if err != nil {
		return nil, err
	}

	for _, a := range attachments {
		input.Attachments = append(input.Attachments, model.AttachmentInput{
			Type:     a.Type,
			URL:      a.URL,
			FileName: a.FileName,
			FileSize: a.FileSize,
		})
	}

	return s.SendMessage(ctx, requesterID, input)
}

func (s *ChatService) findChatReceiver(ctx context.Context, chatID, senderID uuid.UUID) (uuid.UUID, error) {
	members, err := s.membersRepo.GetMembers(ctx, chatID)
	if err != nil {
		return uuid.Nil, err
	}

	for _, m := range members {
		if m.UserID != senderID {
			return m.UserID, nil
		}
	}

	return uuid.Nil, errors.ErrNoMembers
}

func messageText(text *string) string {
	if text == nil {
		return ""
	}

	value := *text
	if len(value) > 100 {
		return value[:100] + "..."
	}

	return value
}
