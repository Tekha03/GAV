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

	existingChatID, err := s.membersRepo.FindPrivateChatBetween(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	if existingChatID != uuid.Nil {
		return s.chatRepo.GetByID(ctx, existingChatID)
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

	if err := s.publishChatCreated(chat.ID, members); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) CreateGroupChat(ctx context.Context, title string, creatorID uuid.UUID, membersIDs []uuid.UUID) (*model.Chat, error) {
	chat := &model.Chat{
		ID:        uuid.New(),
		IsGroup:   true,
		CreatedAt: time.Now(),
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	seenUsers := map[uuid.UUID]struct{}{creatorID: {}}
	members := []*model.ChatMember{{
		ChatID:   chat.ID,
		UserID:   creatorID,
		JoinedAt: time.Now(),
		Role:     model.Admin,
	}}

	for _, id := range membersIDs {
		if _, exists := seenUsers[id]; exists {
			continue
		}
		seenUsers[id] = struct{}{}

		members = append(members, &model.ChatMember{
			ChatID:   chat.ID,
			UserID:   id,
			JoinedAt: time.Now(),
			Role:     model.Member,
		})
	}

	for _, member := range members {
		if err := s.membersRepo.AddMember(ctx, member); err != nil {
			return nil, err
		}
	}

	if err := s.publishChatCreated(chat.ID, members); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) GetChatByID(ctx context.Context, chatID, requesterID uuid.UUID) (*model.Chat, error) {
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

	return chat, nil
}

func (s *ChatService) AddMember(ctx context.Context, userID, chatID, requesterID uuid.UUID) error {
	if err := s.requireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	member := &model.ChatMember{
		ChatID:   chatID,
		UserID:   userID,
		JoinedAt: time.Now(),
		Role:     "member",
	}

	if err := s.membersRepo.AddMember(ctx, member); err != nil {
		return err
	}

	payload, err := json.Marshal(events.ChatMemberAddedData{
		ChatID: chatID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeChatMemberAdded,
		Timestamp: time.Now(),
		Data:      payload,
	}

	return s.publishEvent(event)
}

func (s *ChatService) RemoveMember(ctx context.Context, userID, chatID, requesterID uuid.UUID) error {
	if err := s.requireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	if err := s.membersRepo.RemoveMember(ctx, userID, chatID); err != nil {
		return err
	}

	payload, err := json.Marshal(events.ChatMemberRemovedData{
		ChatID: chatID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeChatMemberRemoved,
		Timestamp: time.Now(),
		Data:      payload,
	}

	return s.publishEvent(event)
}

func (s *ChatService) GetChatMembers(ctx context.Context, chatID, requesterID uuid.UUID) ([]*model.ChatMember, error) {
	if err := s.requireChatMember(ctx, chatID, requesterID); err != nil {
		return nil, err
	}

	members, err := s.membersRepo.GetMembers(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if members == nil {
		return nil, errors.ErrNoMembers
	}

	return members, nil
}

func (s *ChatService) LeaveChat(ctx context.Context, userID, chatID, requesterID uuid.UUID) error {
	if userID != requesterID {
		return errors.ErrChatAccessDenied
	}
	if err := s.requireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	err := s.membersRepo.RemoveMember(ctx, userID, chatID)
	return err
}

func (s *ChatService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*model.Chat, error) {
	chatIDs, err := s.membersRepo.GetUserChats(ctx, userID)
	if err != nil {
		return nil, err
	}
	if chatIDs == nil {
		return nil, errors.ErrNoChats
	}

	var chats []*model.Chat
	for _, id := range chatIDs {
		chat, err := s.chatRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (s *ChatService) UpdateChatTitle(ctx context.Context, chatID, requesterID uuid.UUID, newTitle string) error {
	if err := s.requireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	if err := s.chatRepo.UpdateTitle(ctx, chatID, newTitle); err != nil {
		return errors.ErrTitleUpdate
	}
	return nil
}

func (s *ChatService) UpdateChatPhoto(ctx context.Context, chatID, requesterID uuid.UUID, newPhotoURL string) error {
	if err := s.requireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	if err := s.chatRepo.UpdatePhoto(ctx, chatID, newPhotoURL); err != nil {
		return errors.ErrPhotoUpdate
	}
	return nil
}

func (s *ChatService) publishChatCreated(chatID uuid.UUID, members []*model.ChatMember) error {
	payload, err := json.Marshal(events.ChatCreatedData{
		ChatID:  chatID,
		Members: getMemberIDs(members),
	})
	if err != nil {
		return err
	}

	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeChatCreated,
		Timestamp: time.Now(),
		Data:      payload,
	}

	return s.publishEvent(event)
}

func getMemberIDs(members []*model.ChatMember) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(members))
	for _, member := range members {
		ids = append(ids, member.UserID)
	}
	return ids
}

func (s *ChatService) requireChatMember(ctx context.Context, chatID, userID uuid.UUID) error {
	ok, err := s.membersRepo.MemberExists(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.ErrChatAccessDenied
	}
	return nil
}
