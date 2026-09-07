package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"messenger/internal/constatnts"
	"messenger/internal/model"
	apperrors "shared/app_errors"
	"shared/events"
	"time"

	"github.com/google/uuid"
)

const (
	maxLimit = 100
	minLimit = 0
	midLimit = 50
)

type realtimeMessageEvent struct {
	Type    string             `json:"type"`
	ChatID  uuid.UUID          `json:"chat_id"`
	Message realtimeMessageDTO `json:"message"`
}

type realtimeMessageDTO struct {
	ID          uuid.UUID               `json:"id"`
	ChatID      uuid.UUID               `json:"chat_id"`
	SenderID    uuid.UUID               `json:"sender_id"`
	Text        *string                 `json:"text"`
	ReplyToID   *uuid.UUID              `json:"reply_to_id"`
	CreatedAt   time.Time               `json:"created_at"`
	EditedAt    *time.Time              `json:"edited_at"`
	Attachments []realtimeAttachmentDTO `json:"attachments"`
}

type realtimeAttachmentDTO struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
}

func (s *ChatService) SendMessage(ctx context.Context, requesterID uuid.UUID, input model.SendMessageInput) (*model.Message, error) {
	chat, err := s.chatRepo.GetByID(ctx, input.ChatID)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, apperrors.New(apperrors.ChatNotFound, "chat not found")
	}

	if input.SenderID != requesterID {
		return nil, apperrors.New(apperrors.ChatAccessDenied, "chat access denied")
	}
	if err := s.RequireChatMember(ctx, input.ChatID, requesterID); err != nil {
		return nil, err
	}

	if input.Text == nil && len(input.Attachments) == 0 {
		return nil, apperrors.New(apperrors.MessageContentRequired, "message content is required")
	}
	if input.Text != nil && len(*input.Text) > constatnts.MaxMessageLength {
		return nil, apperrors.New(apperrors.MessageTextTooLong, "message text is too long")
	}
	if len(input.Attachments) > constatnts.MaxAttachments {
		return nil, apperrors.New(apperrors.MessageAttachmentsLimitExceeded, "message attachments limit exceeded")
	}
	if input.ReplyToID != nil {
		msg, err := s.messageRepo.GetByID(ctx, *input.ReplyToID)
		if err != nil {
			return nil, err
		}
		if msg == nil || msg.ChatID != input.ChatID {
			return nil, apperrors.New(apperrors.MessageReplyInvalid, "reply message belongs to another chat")
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

	eventData, err := json.Marshal(events.MessageSentData{
		MessageID: message.ID,
		ChatID:    message.ChatID,
		SenderID:  message.SenderID,
		Text:      messageText(message.Text),
	})
	if err != nil {
		return nil, apperrors.Wrap(apperrors.Internal, "failed to encode message sent data", err)
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeMessageSent,
		Timestamp: time.Now(),
		Data:      eventData,
	}

	eventPayload, err := json.Marshal(event)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.Internal, "failed to encode message sent event", err)
	}

	saveMessage := func(txCtx context.Context) error {
		messageID, err := s.messageRepo.Create(
			txCtx,
			message,
		)
		if err != nil {
			return err
		}

		message.ID = messageID

		attachments := make(
			[]model.Attachment,
			0,
			len(input.Attachments),
		)

		for _, inputAttachment := range input.Attachments {
			attachment := model.Attachment{
				ID:        uuid.New(),
				MessageID: messageID,
				URL:       inputAttachment.URL,
				Type:      inputAttachment.Type,
				FileName:  inputAttachment.FileName,
				FileSize:  inputAttachment.FileSize,
			}

			attachments = append(
				attachments,
				attachment,
			)
		}

		if err := s.attachmentRepo.CreateBatch(
			txCtx,
			attachments,
		); err != nil {
			return err
		}

		message.Attachments = attachments

		if s.outboxRepo != nil {
			if err := s.outboxRepo.Create(txCtx, &model.OutboxEvent{
				ID:            event.EventID,
				EventType:     string(event.EventType),
				Payload:       eventPayload,
				CreatedAt:     event.Timestamp,
				NextAttemptAt: event.Timestamp,
			}); err != nil {
				return err
			}
		}

		return nil
	}

	if s.transactionManager != nil {
		if err := s.transactionManager.WithinTransaction(
			ctx,
			saveMessage,
		); err != nil {
			return nil, err
		}
	} else {
		if err := saveMessage(ctx); err != nil {
			return nil, err
		}
	}

	s.sendRealtimeMessageSent(message)

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
				if err := s.notClient.SendNewMessage(
					context.Background(),
					receiverID,
					senderName,
					text,
					input.ChatID.String(),
				); err != nil {
					slog.Error("failed to send new message notification", "error", err, "chat_id", input.ChatID)
				}
			}()
		}
	}

	if s.outboxRepo == nil {
		if err := s.publishEvent(event); err != nil {
			return nil, err
		}
	}

	return message, nil
}

func (s *ChatService) sendRealtimeMessageSent(message *model.Message) {
	if s.realtime == nil || message == nil {
		return
	}

	payload, err := json.Marshal(realtimeMessageEvent{
		Type:    "message.sent",
		ChatID:  message.ChatID,
		Message: toRealtimeMessageDTO(message),
	})
	if err != nil {
		slog.Error("failed to encode realtime message event", "error", err, "chat_id", message.ChatID)
		return
	}

	s.realtime.SendToChat(message.ChatID, payload)
}

func toRealtimeMessageDTO(message *model.Message) realtimeMessageDTO {
	attachments := make([]realtimeAttachmentDTO, 0, len(message.Attachments))
	for _, attachment := range message.Attachments {
		attachments = append(attachments, realtimeAttachmentDTO{
			ID:        attachment.ID,
			MessageID: attachment.MessageID,
			URL:       attachment.URL,
			Type:      string(attachment.Type),
			FileName:  attachment.FileName,
			FileSize:  attachment.FileSize,
		})
	}

	return realtimeMessageDTO{
		ID:          message.ID,
		ChatID:      message.ChatID,
		SenderID:    message.SenderID,
		Text:        message.Text,
		ReplyToID:   message.ReplyToID,
		CreatedAt:   message.CreatedAt,
		EditedAt:    message.EditedAt,
		Attachments: attachments,
	}
}

func (s *ChatService) EditMessage(ctx context.Context, requesterID, messageID uuid.UUID, newText string) (*model.Message, error) {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, apperrors.New(apperrors.MessageNotFound, "message not found")
	}
	if message.SenderID != requesterID {
		return nil, apperrors.New(apperrors.ChatAccessDenied, "chat access denied")
	}
	if err := s.RequireChatMember(ctx, message.ChatID, requesterID); err != nil {
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
		return nil, apperrors.New(apperrors.MessageNotFound, "message not found")
	}

	payload, err := json.Marshal(events.MessageEditedData{
		MessageID: message.ID,
		ChatID:    message.ChatID,
		Text:      messageText(message.Text),
	})
	if err != nil {
		return nil, apperrors.Wrap(apperrors.Internal, "failed to encode message edited event", err)
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
		return apperrors.New(apperrors.MessageNotFound, "message not found")
	}
	if message.SenderID != requesterID {
		return apperrors.New(apperrors.ChatAccessDenied, "chat access denied")
	}
	if err := s.RequireChatMember(ctx, message.ChatID, requesterID); err != nil {
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
		return apperrors.Wrap(apperrors.Internal, "failed to encode message deleted event", err)
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
	if limit <= minLimit {
		limit = midLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, apperrors.New(apperrors.ChatNotFound, "chat not found")
	}
	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return nil, err
	}

	messages, err := s.messageRepo.GetByChatID(ctx, chatID, limit, cursorID)
	if err != nil {
		return nil, err
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
		return apperrors.New(apperrors.ChatNotFound, "chat not found")
	}
	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	return s.messageRepo.UpdateLastReadMessageForChat(ctx, chatID, requesterID)
}

func (s *ChatService) ForwardMessage(ctx context.Context, requesterID, messageID, targetChatID uuid.UUID) (*model.Message, error) {
	origMsg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if origMsg == nil {
		return nil, apperrors.New(apperrors.MessageNotFound, "message not found")
	}

	if err := s.RequireChatMember(ctx, origMsg.ChatID, requesterID); err != nil {
		return nil, err
	}
	if err := s.RequireChatMember(ctx, targetChatID, requesterID); err != nil {
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

	return uuid.Nil, apperrors.New(apperrors.ChatMemberNotFound, "chat receiver not found")
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
