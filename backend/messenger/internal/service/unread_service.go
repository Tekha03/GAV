package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *ChatService) GetChatUnreadCount(ctx context.Context, chatID, userID uuid.UUID) (int, error) {
	if err := s.requireChatMember(ctx, chatID, userID); err != nil {
		return 0, err
	}

	lastReadID, err := s.membersRepo.GetLastReadMessageID(ctx, chatID, userID)
	if err != nil {
		return 0, err
	}

	messages, err := s.messageRepo.GetByChatID(ctx, chatID, 0, nil)
	if err != nil {
		return 0, err
	}

	lastReadTime := time.Time{}

	if lastReadID != uuid.Nil {
		for _, msg := range messages {
			if msg.ID == lastReadID {
				lastReadTime = msg.CreatedAt
				break
			}
		}
	}

	count := 0
	for _, msg := range messages {
		if msg.CreatedAt.After(lastReadTime) {
			count++
		}
	}

	return count, nil
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
