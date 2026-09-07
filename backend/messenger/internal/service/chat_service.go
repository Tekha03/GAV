package service

import (
	"context"
	"encoding/json"
	"messenger/internal/model"
	apperrors "shared/app_errors"
	"shared/events"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *ChatService) CreatePrivateChat(ctx context.Context, userID1, userID2 uuid.UUID) (*model.Chat, error) {
	if userID1 == userID2 {
		return nil, apperrors.New(apperrors.ChatSelfNotAllowed, "cannot create a chat with yourself")
	}

	privateKey := privateChatKey(userID1, userID2)
	var result *model.Chat
	var createdMembers []*model.ChatMember

	create := func(txCtx context.Context) error {
		existingChat, err := s.chatRepo.GetByPrivateKey(txCtx, privateKey)
		if err != nil {
			return err
		}
		if existingChat != nil {
			result = existingChat
			return nil
		}

		// Support private chats created before private_key was introduced.
		existingChatID, err := s.membersRepo.FindPrivateChatBetween(txCtx, userID1, userID2)
		if err != nil {
			return err
		}
		if existingChatID != uuid.Nil {
			result, err = s.chatRepo.GetByID(txCtx, existingChatID)
			return err
		}

		chat := &model.Chat{
			ID:         uuid.New(),
			IsGroup:    false,
			PrivateKey: &privateKey,
			CreatedAt:  time.Now(),
		}
		if err := s.chatRepo.Create(txCtx, chat); err != nil {
			return err
		}

		members := []*model.ChatMember{
			{ChatID: chat.ID, UserID: userID1, JoinedAt: time.Now(), Role: model.Member},
			{ChatID: chat.ID, UserID: userID2, JoinedAt: time.Now(), Role: model.Member},
		}
		for _, member := range members {
			if err := s.membersRepo.AddMember(txCtx, member); err != nil {
				return err
			}
		}

		result = chat
		createdMembers = members
		return nil
	}

	var err error
	if s.transactionManager != nil {
		err = s.transactionManager.WithPrivateChatLock(ctx, privateKey, create)
	} else {
		err = create(ctx)
	}
	if err != nil {
		return nil, err
	}

	if len(createdMembers) > 0 {
		if err := s.publishChatCreated(result.ID, createdMembers); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func privateChatKey(userID1, userID2 uuid.UUID) string {
	ids := []string{userID1.String(), userID2.String()}
	sort.Strings(ids)
	return strings.Join(ids, ":")
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
		return nil, apperrors.New(apperrors.ChatNotFound, "chat not found")
	}

	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) AddMember(ctx context.Context, chatID, userID, requesterID uuid.UUID) error {
	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	role, err := s.membersRepo.GetRole(ctx, requesterID, chatID)
	if err != nil {
		return err
	}

	if *role != model.Admin {
		return apperrors.New(apperrors.ChatAccessDenied, "chat access denied")
	}

	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.IsGroup != true {
		return apperrors.New(apperrors.ChatGroupRequired, "operation is available only for group chats")
	}

	member := &model.ChatMember{
		ChatID:   chatID,
		UserID:   userID,
		JoinedAt: time.Now(),
		Role:     model.Member,
	}

	if err := s.membersRepo.AddMember(ctx, member); err != nil {
		return err
	}

	payload, err := json.Marshal(events.ChatMemberAddedData{
		ChatID: chatID,
		UserID: userID,
	})
	if err != nil {
		return apperrors.Wrap(apperrors.Internal, "failed to encode chat member added event", err)
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
	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	role, err := s.membersRepo.GetRole(ctx, requesterID, chatID)
	if err != nil {
		return err
	}

	if *role != model.Admin {
		return apperrors.New(apperrors.ChatAccessDenied, "chat access denied")
	}

	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.IsGroup != true {
		return apperrors.New(apperrors.ChatGroupRequired, "operation is available only for group chats")
	}

	if err := s.membersRepo.RemoveMember(ctx, userID, chatID); err != nil {
		return err
	}

	payload, err := json.Marshal(events.ChatMemberRemovedData{
		ChatID: chatID,
		UserID: userID,
	})
	if err != nil {
		return apperrors.Wrap(apperrors.Internal, "failed to encode chat member removed event", err)
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
	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return nil, err
	}

	members, err := s.membersRepo.GetMembers(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (s *ChatService) LeaveChat(ctx context.Context, chatID, requesterID uuid.UUID) error {

	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	err := s.membersRepo.RemoveMember(ctx, requesterID, chatID)
	return err
}

func (s *ChatService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*model.Chat, error) {
	chatIDs, err := s.membersRepo.GetUserChats(ctx, userID)
	if err != nil {
		return nil, err
	}
	chats := make([]*model.Chat, 0, len(chatIDs))
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
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.IsGroup != true {
		return apperrors.New(apperrors.ChatGroupRequired, "operation is available only for group chats")
	}

	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	if err := s.chatRepo.UpdateTitle(ctx, chatID, newTitle); err != nil {
		return err
	}
	return nil
}

func (s *ChatService) UpdateChatPhoto(ctx context.Context, chatID, requesterID uuid.UUID, newPhotoURL string) error {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.IsGroup != true {
		return apperrors.New(apperrors.ChatGroupRequired, "operation is available only for group chats")
	}

	if err := s.RequireChatMember(ctx, chatID, requesterID); err != nil {
		return err
	}

	if err := s.chatRepo.UpdatePhoto(ctx, chatID, newPhotoURL); err != nil {
		return err
	}
	return nil
}

func (s *ChatService) publishChatCreated(chatID uuid.UUID, members []*model.ChatMember) error {
	payload, err := json.Marshal(events.ChatCreatedData{
		ChatID:  chatID,
		Members: getMemberIDs(members),
	})
	if err != nil {
		return apperrors.Wrap(apperrors.Internal, "failed to encode chat created event", err)
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

func (s *ChatService) RequireChatMember(ctx context.Context, chatID, userID uuid.UUID) error {
	ok, err := s.membersRepo.MemberExists(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return apperrors.New(apperrors.ChatAccessDenied, "chat access denied")
	}
	return nil
}
