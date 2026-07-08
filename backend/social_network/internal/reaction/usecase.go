package reaction

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

func (u *UseCase) OnReactionAdded(ctx context.Context, data events.ReactionAddedData) {
	u.store.AddReaction(data.MessageID, data.UserID, data.Reaction)

	if u.notifier == nil {
		return
	}

	message, ok := u.store.Message(data.MessageID)
	if !ok || message.SenderID == data.UserID {
		return
	}

	if err := u.notifier.NotifyMessageReaction(ctx, message.SenderID, data.UserID, data.MessageID); err != nil {
		log.Printf("reaction added notification error: %v", err)
	}
}

func (u *UseCase) OnReactionRemoved(ctx context.Context, data events.ReactionRemovedData) {
	u.store.RemoveReaction(data.MessageID, data.UserID)
}
