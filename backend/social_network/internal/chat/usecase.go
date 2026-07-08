package chat

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

func (u *UseCase) OnChatCreated(ctx context.Context, data events.ChatCreatedData) {
	u.store.CreateChat(data.ChatID, data.Members)

	if u.notifier == nil {
		return
	}

	for _, memberID := range data.Members {
		if err := u.notifier.NotifyChatInvite(ctx, memberID, data.ChatID); err != nil {
			log.Printf("chat created notification error: %v", err)
		}
	}
}

func (u *UseCase) OnChatMemberAdded(ctx context.Context, data events.ChatMemberAddedData) {
	u.store.AddMember(data.ChatID, data.UserID)

	if u.notifier == nil {
		return
	}

	if err := u.notifier.NotifyChatInvite(ctx, data.UserID, data.ChatID); err != nil {
		log.Printf("chat member added notification error: %v", err)
	}

	log.Println("handle chat member added:", data.ChatID, data.UserID)
}

func (u *UseCase) OnChatMemberRemoved(ctx context.Context, data events.ChatMemberRemovedData) {
	u.store.RemoveMember(data.ChatID, data.UserID)
}

func (u *UseCase) OnChatDeleted(ctx context.Context, data events.ChatDeletedData) {
	u.store.DeleteChat(data.ChatID)
}
