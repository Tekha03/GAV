package chat

import (
	"context"
	"testing"

	"shared/events"
	messengersync "social_network/internal/messenger_sync"

	"github.com/google/uuid"
)

func TestUseCaseChatLifecycleWithoutNotifier(t *testing.T) {
	store := messengersync.NewStore()
	useCase := NewUseCase(store, nil)
	chatID, firstUser, secondUser := uuid.New(), uuid.New(), uuid.New()

	useCase.OnChatCreated(context.Background(), events.ChatCreatedData{ChatID: chatID, Members: []uuid.UUID{firstUser}})
	useCase.OnChatMemberAdded(context.Background(), events.ChatMemberAddedData{ChatID: chatID, UserID: secondUser})
	if members := store.ChatMembers(chatID); len(members) != 2 {
		t.Fatalf("members after add = %+v", members)
	}
	useCase.OnChatMemberRemoved(context.Background(), events.ChatMemberRemovedData{ChatID: chatID, UserID: firstUser})
	if members := store.ChatMembers(chatID); len(members) != 1 || members[0] != secondUser {
		t.Fatalf("members after remove = %+v", members)
	}
	useCase.OnChatDeleted(context.Background(), events.ChatDeletedData{ChatID: chatID})
}
