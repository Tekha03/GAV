package message

import (
	"context"
	"log"
	"shared/events"
	messengersync "social_network/internal/messenger_sync"
	"social_network/internal/notification"
)

type UseCase struct {
	store    *messengersync.Store
	notifier notification.NotificationService
}

func NewUseCase(store *messengersync.Store, notifier notification.NotificationService) *UseCase {
	return &UseCase{
		store:    store,
		notifier: notifier,
	}
}

func (u *UseCase) OnMessageSent(ctx context.Context, data events.MessageSentData) {
	u.store.SaveMessage(data.MessageID, data.ChatID, data.SenderID, data.Text)

	if u.notifier == nil {
		return
	}

	memberIDs := u.store.ChatMembers(data.ChatID)
	for _, memberID := range memberIDs {
		if memberID == data.SenderID {
			continue
		}
		if err := u.notifier.NotifyNewMessage(ctx, memberID, data.SenderID, data.MessageID); err != nil {
			log.Printf("message sent notification error: %v", err)
		}
	}
}

func (u *UseCase) OnMessageEdited(ctx context.Context, data events.MessageEditedData) {
	u.store.EditMessage(data.MessageID, data.Text)
}

func (u *UseCase) OnMessageDeleted(ctx context.Context, data events.MessageDeletedData) {
	u.store.DeleteMessage(data.MessageID)
}
