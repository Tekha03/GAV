package messengersync

import (
	"testing"

	"github.com/google/uuid"
)

func TestStoreChatLifecycle(t *testing.T) {
	store := NewStore()
	chatID, firstUser, secondUser := uuid.New(), uuid.New(), uuid.New()

	store.CreateChat(chatID, []uuid.UUID{firstUser})
	if members := store.ChatMembers(chatID); len(members) != 1 || members[0] != firstUser {
		t.Fatalf("ChatMembers() = %+v", members)
	}
	store.AddMember(chatID, secondUser)
	if members := store.ChatMembers(chatID); len(members) != 2 {
		t.Fatalf("members after add = %+v", members)
	}
	store.RemoveMember(chatID, firstUser)
	if members := store.ChatMembers(chatID); len(members) != 1 || members[0] != secondUser {
		t.Fatalf("members after remove = %+v", members)
	}
	store.DeleteChat(chatID)
	if !store.chats[chatID].Deleted {
		t.Fatal("DeleteChat() did not mark chat deleted")
	}
}

func TestStoreToleratesOutOfOrderChatEvents(t *testing.T) {
	store := NewStore()
	chatID, userID := uuid.New(), uuid.New()
	store.AddMember(chatID, userID)
	if members := store.ChatMembers(chatID); len(members) != 1 || members[0] != userID {
		t.Fatalf("ChatMembers() = %+v", members)
	}
	store.RemoveMember(uuid.New(), userID)
	store.DeleteChat(uuid.New())
	if members := store.ChatMembers(uuid.New()); members != nil {
		t.Fatalf("missing chat members = %+v", members)
	}
}

func TestStoreMessageAndReactionLifecycle(t *testing.T) {
	store := NewStore()
	messageID, chatID, senderID, reactorID := uuid.New(), uuid.New(), uuid.New(), uuid.New()

	store.SaveMessage(messageID, chatID, senderID, "hello")
	message, ok := store.Message(messageID)
	if !ok || message.Text != "hello" || message.SenderID != senderID {
		t.Fatalf("Message() = %+v, %v", message, ok)
	}
	message.Text = "external mutation"
	stored, _ := store.Message(messageID)
	if stored.Text != "hello" {
		t.Fatal("Message() returned mutable store state")
	}

	store.EditMessage(messageID, "edited")
	store.AddReaction(messageID, reactorID, "👍")
	if stored, _ = store.Message(messageID); stored.Text != "edited" {
		t.Fatalf("edited message = %+v", stored)
	}
	if store.reactions[messageID][reactorID] != "👍" {
		t.Fatalf("reaction = %q", store.reactions[messageID][reactorID])
	}
	store.RemoveReaction(messageID, reactorID)
	if _, exists := store.reactions[messageID][reactorID]; exists {
		t.Fatal("reaction was not removed")
	}
	store.DeleteMessage(messageID)
	if stored, _ = store.Message(messageID); !stored.Deleted {
		t.Fatalf("deleted message = %+v", stored)
	}
}

func TestStoreToleratesMissingMessageEvents(t *testing.T) {
	store := NewStore()
	messageID, userID := uuid.New(), uuid.New()
	store.EditMessage(messageID, "ignored")
	store.DeleteMessage(messageID)
	store.RemoveReaction(messageID, userID)
	if message, ok := store.Message(messageID); ok || message != nil {
		t.Fatalf("Message() = %+v, %v", message, ok)
	}
}
