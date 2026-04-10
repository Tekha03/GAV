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

func (s *ChatService) CreatePrivateChat(ctx context.Context, userID1, userID2 uuid.UUID) (*model.Chat, error) {
	if userID1 == userID2 {
		return nil, errors.ErrChatWithYourSelf
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

func (s *ChatService) CreateGroupChat(ctx context.Context, title string, creatorID uuid.UUID, membersIDs []uuid.UUID) (*model.Chat, error) {
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

    payload, _ := json.Marshal(events.ChatCreatedData{
        ChatID: chat.ID,
        Members: getMemberIDs(members),
    })

    event := events.Event{
        EventID:    uuid.New(),
        EventType:  events.EventTypeChatCreated,
        Timestamp:  time.Now(),
        Data:       payload,
    }

    if err := s.producer.PublishEvent(event); err != nil {
        return nil, err
    }

	return chat, nil
}

func (s *ChatService) GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
    chat, err := s.chatRepo.GetByID(ctx, chatID)
    if err != nil {
        return nil, err
    }

    if chat == nil {
        return nil, errors.ErrChatNotFound
    }

    return chat, nil
}

func (s *ChatService) AddMember(ctx context.Context, userID, chatID uuid.UUID) error {
	member := &model.ChatMember{
		ChatID:   chatID,
		UserID:   userID,
		JoinedAt: time.Now(),
		Role:     "member",
	}

    payload, _ := json.Marshal(events.ChatMemberAddedData{
        ChatID: chatID,
        UserID: userID,
    })

    event := events.Event{
        EventID:    uuid.New(),
        EventType:  events.EventTypeChatMemberAdded,
        Timestamp:  time.Now(),
        Data:       payload,
    }

    if err := s.producer.PublishEvent(event); err != nil {
        return err
    }

	return s.membersRepo.AddMember(ctx, member)
}

func (s *ChatService) RemoveMember(ctx context.Context, userID, chatID uuid.UUID) error {
    payload, _ := json.Marshal(events.ChatMemberRemovedData{
        ChatID: chatID,
        UserID: userID,
    })

    event := events.Event{
        EventID:    uuid.New(),
        EventType:  events.EventTypeChatMemberRemoved,
        Timestamp:  time.Now(),
        Data:       payload,
    }

    if err := s.producer.PublishEvent(event); err != nil {
        return err
    }

    return s.membersRepo.RemoveMember(ctx, userID, chatID)
}

func (s *ChatService) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error) {
    members, err := s.membersRepo.GetMembers(ctx, chatID)
    if err != nil {
        return nil, err
    }

    if members == nil {
        return nil, errors.ErrNoMembers
    }

    return members, nil
}

func (s *ChatService) LeaveChat(ctx context.Context, userID, chatID uuid.UUID) error {
    err := s.membersRepo.RemoveMember(ctx, userID, chatID)
    return err
}

func (s *ChatService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*model.Chat, error) {
    chat_ids, err := s.membersRepo.GetUserChats(ctx, userID)
    if err != nil {
        return nil, err
    }
    if chat_ids == nil {
        return nil, errors.ErrNoChats
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

func (s *ChatService) UpdateChatTitle(ctx context.Context, chatID uuid.UUID, newTitle string) error {
    err := s.chatRepo.UpdateTitle(ctx, chatID, newTitle)
    if err != nil {
        return errors.ErrTitleUpdate
    }
    return nil
}

func (s *ChatService) UpdateChatPhoto(ctx context.Context, chatID uuid.UUID, newPhotoURL string) error {
    err := s.chatRepo.UpdatePhoto(ctx, chatID, newPhotoURL)
    if err != nil {
        return errors.ErrPhotoUpdate
    }

    return nil
}

func getMemberIDs(members []*model.ChatMember) []uuid.UUID {
    var ids []uuid.UUID
    for _, m := range members {
        ids = append(ids, m.UserID)
    }
    return ids
}
