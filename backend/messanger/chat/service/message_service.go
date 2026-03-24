package service

import (
	"context"
	"messanger/chat/constatnts"
	"messanger/chat/errors"
	"messanger/chat/model"
	"time"

	"github.com/google/uuid"
)

func (s *ChatService) SendMessage(ctx context.Context, input model.SendMessageInput) (*model.Message, error) {
    chat, err := s.chatRepo.GetByID(ctx, input.ChatID)
    if err != nil {
        return nil, err
    }

    if chat == nil {
        return nil, errors.ErrChatNotFound
    }

    members, err := s.membersRepo.GetMembers(ctx, input.ChatID)
    if err != nil {
        return nil, err
    }

    if members == nil {
        return nil, errors.ErrNoMembers
    }

    if input.Text == nil && len(input.Attachments) == 0 {
        return nil, errors.ErrEmptyMessage
    }

    if input.ReplyToID != nil {
        msg, err := s.messageRepo.GetByID(ctx, *input.ReplyToID)
        if err != nil {
            return nil, err
        }
    
        if msg.ChatID != input.ChatID {
            return nil, errors.ErrInvalidReply
        }
    }

	if len(*input.Text) > constatnts.MaxMessageLength {
		return nil, errors.ErrTextOverLength
	}

    message := &model.Message{
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

	if len(input.Attachments) > constatnts.MaxAttachments {
		return nil, errors.ErrAttachmentsOverLength
	}

    for _, att :=range input.Attachments {
        attachment := model.Attachment {
            MessageID: msgID,
            URL: att.URL,
            Type: att.Type,
            FileName: att.FileName,
            FileSize: att.FileSize,
        }

        if err := s.attachmentRepo.Create(ctx, &attachment); err != nil {
            return nil, err
        }
    }

    return message, nil
}

func (s *ChatService) EditMessage(ctx context.Context, messageID uuid.UUID, newText string) (*model.Message, error) {
    message, err := s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return nil, err
    }

    if message == nil {
        return nil, errors.ErrMessageNotFound
    }

    err = s.messageRepo.UpdateText(ctx, messageID, newText)
    if err != nil {
        return nil, err
    }

    message, err = s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return nil, err
    }

    if message == nil {
        return nil, errors.ErrMessageNotFound
    }

    return message, nil

}
func (s *ChatService) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
    message, err := s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return err
    }

    if message == nil {
        return errors.ErrMessageNotFound
    }

    err = s.messageRepo.Delete(ctx, messageID)
    return err
}
func (s *ChatService) GetMessages(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
    chat, err := s.chatRepo.GetByID(ctx, chatID)

    if err != nil {
        return nil, err
    }
    
    if chat == nil {
        return nil, errors.ErrChatNotFound
    }

    messages, err := s.messageRepo.GetByChatID(ctx, chatID, limit, cursorID)

    if err != nil {
        return nil, err
    }

    if messages == nil {
        return nil, errors.ErrMessageNotFound
    }

    return messages, nil
}
func (s *ChatService) MarkAsRead(ctx context.Context, chatID, userID uuid.UUID) error {
    chat, err := s.chatRepo.GetByID(ctx, chatID)

    if err != nil {
        return err
    }

    if chat == nil {
        return errors.ErrChatNotFound
    }

    err = s.messageRepo.UpdateReadAtForChat(ctx, chatID, userID, time.Now())
    if err != nil {
        return err
    }

    return nil
}

func (s *ChatService) ForwardMessage(ctx context.Context, messageID, targetChatID, senderID uuid.UUID) (*model.Message, error) {
    origMsg, err := s.messageRepo.GetByID(ctx, messageID)
    if err != nil {
        return nil, err
    }
    if origMsg == nil {
        return nil, errors.ErrMessageNotFound
    }

	members, err := s.membersRepo.GetMembers(ctx, targetChatID)
    if err != nil {
        return nil, err
    }

	isMember := false
	for _, member := range members {
		if member.UserID == senderID {
			isMember = true
			break
		}
	}

	if !isMember {
		return nil, errors.ErrIsNotMember
	}

    input := model.SendMessageInput{
        ChatID:      targetChatID,
        SenderID:    senderID,
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

    newMsg, err := s.SendMessage(ctx, input)
    if err != nil {
        return nil, err
    }

    return newMsg, nil
}