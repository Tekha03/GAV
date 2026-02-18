package chat

import (
	"context"
	"errors"
	"gav/internal/chat/model"
	"gav/internal/chat/repository"
	"time"

	"github.com/google/uuid"
)

var (
    ErrChatWithYourSelf = errors.New("can not create private chat with yourself of add yourself to group")
    ErrChatNotFound = errors.New("chat not found")
    ErrNoMembers = errors.New("no members in chat")
    ErrNoChats = errors.New("no chats")
    ErrTitleUpdate = errors.New("error to update chat title")
    ErrPhotoUpdate = errors.New("error to update chat photo")
    ErrEmptyMessage = errors.New("message content can not be empty")
    ErrMessageNotFound = errors.New("message not found")
)

type service struct {
	chatRepo      repository.ChatRepository
	membersRepo   repository.ChatMemberRepository
	messageRepo   repository.MessageRepository
	attachmenRepo repository.AttachmentRepository
	reactionRepo  repository.ReactionRepository
}

func NewService(
	chatRepo repository.ChatRepository, membersRepo repository.ChatMemberRepository,
	messageRepo repository.MessageRepository,
	reactionRepo repository.ReactionRepository,
) Service {
	return &service{
		chatRepo:     chatRepo,
		membersRepo:  membersRepo,
		messageRepo:  messageRepo,
		reactionRepo: reactionRepo,
	}
}

func (s *service) CreatePrivateChat(ctx context.Context, userID1, userID2 uuid.UUID) (*model.Chat, error) {
	if userID1 == userID2 {
		return nil, ErrChatWithYourSelf
	}

	chat := &model.Chat{
		ID:        uuid.New(),
		IsGroup:   false,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	members := []*model.ChatMember{
		{ChatID: chat.ID, UserID: userID1, JoinedAt: time.Now(), Role: model.Member},
		{ChatID: chat.ID, UserID: userID2, JoinedAt: time.Now(), Role: model.Member},
	}

	for _, m := range members {
		if err := s.membersRepo.AddMember(ctx, m); err != nil {
			return nil, err
		}
	}

	return chat, nil
}

func (s *service) CreateGroupChat(ctx context.Context, title string, creatorID uuid.UUID, membersIDs []uuid.UUID) (*model.Chat, error) {

    chat := &model.Chat{
        ID: uuid.New(),
        IsGroup: true,
        CreatedAt: time.Now(),
    }
    
    var members []*model.ChatMember
    members = append(members, &model.ChatMember{
        ChatID: chat.ID,
        UserID: creatorID,
        JoinedAt: time.Now(),
        Role: model.Admin,
    })

    for _, id := range membersIDs {
        if creatorID == id {
            continue
        }
        member := &model.ChatMember{
            ChatID: chat.ID,
            UserID: id,
            JoinedAt: time.Now(),
            Role: model.Member,
        }

        members = append(members, member)
    }

    for _, m := range members {
		if err := s.membersRepo.AddMember(ctx, m); err != nil {
			return nil, err
		}
	}

	return chat, nil
}

func (s *service) GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
    chat, err := s.chatRepo.GetByID(ctx, chatID)
    if err != nil {
        return nil, err
    } 

    if chat == nil {
        return nil, ErrChatNotFound
    }

    return chat, nil
}

func (s *service) AddMember(ctx context.Context, userID, chatID uuid.UUID) error {
	member := &model.ChatMember{
		ChatID:   chatID,
		UserID:   userID,
		JoinedAt: time.Now(),
		Role:     "member",
	}

	return s.membersRepo.AddMember(ctx, member)
}

func (s *service) RemoveMember(ctx context.Context, userID, chatID uuid.UUID) error {
    return s.membersRepo.RemoveMember(ctx, userID, chatID)
}

func (s *service) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error) {
    members, err := s.membersRepo.GetMembers(ctx, chatID)
    if err != nil {
        return nil, err
    }

    if members == nil {
        return nil, ErrNoMembers
    }

    return members, nil
}

func (s *service) LeaveChat(ctx context.Context, userID, chatID uuid.UUID) error {
    err := s.membersRepo.RemoveMember(ctx, userID, chatID)
    return err
}

func (s *service) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*model.Chat, error) {
    chat_ids, err := s.membersRepo.GetUserChats(ctx, userID)
    if err != nil {
        return nil, err
    }
    if chat_ids == nil {
        return nil, ErrNoChats
    }

    var chats []*model.Chat
    for _, id := range chat_ids {
        chat, err := s.chatRepo.GetByID(ctx, id)
        if err != nil {
            return nil, err
        }

        chats = append(chats, chat)
    }

    return chats, nil
}

func (s *service) UpdateChatTitle(ctx context.Context, chatID uuid.UUID, newTitle string) error {
    err := s.chatRepo.UpdateTitle(ctx, chatID, newTitle)
    if err != nil {
        return ErrTitleUpdate
    }
    return nil
}

func (s *service) UpdateChatPhoto(ctx context.Context, chatID uuid.UUID, newPhotoURL string) error {
    err := s.chatRepo.UpdatePhoto(ctx, chatID, newPhotoURL)
    if err != nil {
        return ErrPhotoUpdate
    }

    return nil
}

func (s *service) SendMessage(ctx context.Context, input model.SendMessageInput) (*model.Message, error) {
    chat, err := s.chatRepo.GetByID(ctx, input.ChatID)
    if err != nil {
        return nil, err
    }

    if chat == nil {
        return nil, ErrChatNotFound
    }

    member, err := s.membersRepo.GetMembers(ctx, input.ChatID)
    if err != nil {
        return nil, err
    }

    if member == nil {
        return nil, ErrNoMembers
    }

    message := &model.Message{
		ID:        uuid.New(),
		ChatID:    input.ChatID,
		SenderID:  input.SenderID,
		Text:      input.Text,
		ReplyToID: input.ReplyToID,
		CreatedAt: time.Now(),
	}

    if err := s.messageRepo.Create(ctx, message); err != nil {
        return nil, err
    }

    return message, nil
}

func (s *service) EditMessage(ctx context.Context, messageID uuid.UUID, newText string) (*model.Message, error) {
    message, err := s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return nil, err
    }

    if message == nil {
        return nil, ErrMessageNotFound
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
        return nil, ErrMessageNotFound
    }

    return message, nil

}
func (s *service) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
    message, err := s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return err
    }

    if message == nil {
        return ErrMessageNotFound
    }

    err = s.messageRepo.Delete(ctx, messageID)
    return err
}
func (s *service) GetMessages(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
    chat, err := s.chatRepo.GetByID(ctx, chatID)

    if err != nil {
        return nil, err
    }
    
    if chat == nil {
        return nil, ErrChatNotFound
    }

    messages, err := s.messageRepo.GetByChatID(ctx, chatID, limit, cursorID)

    if err != nil {
        return nil, err
    }

    if messages == nil {
        return nil, ErrMessageNotFound
    }

    return messages, nil
}
func (s *service) MarkAsRead(ctx context.Context, chatID, userID uuid.UUID) error {
    chat, err := s.chatRepo.GetByID(ctx, chatID)

    if err != nil {
        return err
    }

    if chat == nil {
        return ErrChatNotFound
    }

    err = s.messageRepo.UpdateReadAtForChat(ctx, chatID, userID, time.Now())
    if err != nil {
        return err
    }

    return nil
}

func (s *service) AddReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
    msg, err := s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return err
    }

    if msg == nil {
        return ErrMessageNotFound
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
func (s *service) RemoveReaction(ctx context.Context, messageID, userID uuid.UUID) error {
    msg, err := s.messageRepo.GetByID(ctx, messageID)

    if err != nil {
        return err
    }

    if msg == nil {
        return ErrMessageNotFound
    }

    err = s.reactionRepo.Remove(ctx, messageID, userID)

    return err
}