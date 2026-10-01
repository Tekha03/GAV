package message

import (
	"context"
	"errors"
	"testing"

	"shared/events"
	messengersync "social_network/internal/messenger_sync"
	"social_network/internal/notification"

	"github.com/google/uuid"
)

type notifierStub struct {
	receivers []uuid.UUID
	sender    uuid.UUID
	entity    uuid.UUID
	err       error
}

func (*notifierStub) NotifyLike(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error { return nil }
func (*notifierStub) NotifyComment(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error {
	return nil
}
func (*notifierStub) NotifyFollow(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (n *notifierStub) NotifyNewMessage(_ context.Context, receiverID, senderID, entityID uuid.UUID) error {
	n.receivers = append(n.receivers, receiverID)
	n.sender, n.entity = senderID, entityID
	return n.err
}
func (*notifierStub) NotifyChatInvite(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (*notifierStub) NotifyMessageReaction(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error {
	return nil
}
func (*notifierStub) GetInAppNotifications(context.Context, uuid.UUID) ([]*notification.Notification, error) {
	return nil, nil
}
func (*notifierStub) MarkInAppNotificationAsRead(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}

func TestUseCaseMessageLifecycleAndNotifications(t *testing.T) {
	store := messengersync.NewStore()
	notifier := &notifierStub{}
	useCase := NewUseCase(store, notifier)
	chatID, senderID, firstReceiver, secondReceiver, messageID := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	store.CreateChat(chatID, []uuid.UUID{senderID, firstReceiver, secondReceiver})

	useCase.OnMessageSent(context.Background(), events.MessageSentData{
		MessageID: messageID,
		ChatID:    chatID,
		SenderID:  senderID,
		Text:      "hello",
	})
	message, ok := store.Message(messageID)
	if !ok || message.Text != "hello" {
		t.Fatalf("stored message = %+v, %v", message, ok)
	}
	if len(notifier.receivers) != 2 || notifier.sender != senderID || notifier.entity != messageID {
		t.Fatalf("notifications = %+v sender=%s entity=%s", notifier.receivers, notifier.sender, notifier.entity)
	}
	for _, receiver := range notifier.receivers {
		if receiver == senderID {
			t.Fatal("sender received its own notification")
		}
	}

	useCase.OnMessageEdited(context.Background(), events.MessageEditedData{MessageID: messageID, Text: "edited"})
	message, _ = store.Message(messageID)
	if message.Text != "edited" {
		t.Fatalf("edited message = %+v", message)
	}
	useCase.OnMessageDeleted(context.Background(), events.MessageDeletedData{MessageID: messageID})
	message, _ = store.Message(messageID)
	if !message.Deleted {
		t.Fatalf("deleted message = %+v", message)
	}
}

func TestUseCaseWorksWithoutNotifierAndIgnoresNotifierErrors(t *testing.T) {
	store := messengersync.NewStore()
	chatID, senderID, receiverID := uuid.New(), uuid.New(), uuid.New()
	store.CreateChat(chatID, []uuid.UUID{senderID, receiverID})

	NewUseCase(store, nil).OnMessageSent(context.Background(), events.MessageSentData{
		MessageID: uuid.New(), ChatID: chatID, SenderID: senderID,
	})
	notifier := &notifierStub{err: errors.New("push failed")}
	NewUseCase(store, notifier).OnMessageSent(context.Background(), events.MessageSentData{
		MessageID: uuid.New(), ChatID: chatID, SenderID: senderID,
	})
	if len(notifier.receivers) != 1 {
		t.Fatalf("notification attempts = %d", len(notifier.receivers))
	}
}
