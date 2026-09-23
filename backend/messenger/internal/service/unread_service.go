package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *ChatService) GetChatUnreadCount(ctx context.Context, chatID, userID uuid.UUID) (int, error) {
	if err := s.RequireChatMember(ctx, chatID, userID); err != nil {
		return 0, err
	}

	lastReadID, err := s.membersRepo.GetLastReadMessageID(ctx, chatID, userID)
	if err != nil {
		return 0, err
	}

	return s.messageRepo.CountUnreadForChat(ctx, chatID, userID, lastReadID)
}

func (s *ChatService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	// 1. Получаем все чаты пользователя
	chatIDs, err := s.membersRepo.GetUserChats(ctx, userID)
	if err != nil {
		return 0, err
	}

	totalUnread := 0
	for _, chatID := range chatIDs {
		unread, err := s.GetChatUnreadCount(ctx, chatID, userID)
		if err != nil {
			return 0, err
		}
		totalUnread += unread
	}

	return totalUnread, nil
}
