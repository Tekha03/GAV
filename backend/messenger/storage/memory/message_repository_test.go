package memory

import (
	"context"
	"testing"
	"time"

	"messenger/internal/model"

	"github.com/google/uuid"
)

func TestMessageRepositoryLifecyclePaginationAndUnreadCount(t *testing.T) {
	ctx := context.Background()
	repo := NewMessageRepository()
	chatID, otherChatID := uuid.New(), uuid.New()
	currentUser, otherUser := uuid.New(), uuid.New()
	baseTime := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)

	messages := []*model.Message{
		{ID: uuid.New(), ChatID: chatID, SenderID: otherUser},
		{ID: uuid.New(), ChatID: chatID, SenderID: currentUser},
		{ID: uuid.New(), ChatID: chatID, SenderID: otherUser},
		{ID: uuid.New(), ChatID: otherChatID, SenderID: otherUser},
	}
	for index, message := range messages {
		if _, err := repo.Create(ctx, message); err != nil {
			t.Fatal(err)
		}
		message.CreatedAt = baseTime.Add(time.Duration(index) * time.Minute)
	}

	if _, err := repo.Create(ctx, messages[0]); err == nil {
		t.Fatal("duplicate Create() error = nil")
	}
	generated := &model.Message{ChatID: chatID, SenderID: otherUser}
	if id, err := repo.Create(ctx, generated); err != nil || id == uuid.Nil || generated.ID != id {
		t.Fatalf("generated Create() id=%s error=%v", id, err)
	}
	generated.CreatedAt = baseTime.Add(4 * time.Minute)

	text := "updated"
	if err := repo.UpdateText(ctx, messages[0].ID, text); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.GetByID(ctx, messages[0].ID)
	if err != nil || loaded.Text == nil || *loaded.Text != text || loaded.EditedAt == nil {
		t.Fatalf("GetByID() = %+v, error=%v", loaded, err)
	}
	if err := repo.UpdateText(ctx, uuid.New(), text); err == nil {
		t.Fatal("UpdateText() missing message error = nil")
	}

	latestID, err := repo.GetLastMessageIDForChat(ctx, chatID)
	if err != nil || latestID != generated.ID {
		t.Fatalf("GetLastMessageIDForChat() = %s, %v", latestID, err)
	}
	page, err := repo.GetByChatID(ctx, chatID, 2, nil)
	if err != nil || len(page) != 2 || page[0].ID != generated.ID || page[1].ID != messages[2].ID {
		t.Fatalf("first page = %+v, error=%v", page, err)
	}
	cursor := page[1].ID
	nextPage, err := repo.GetByChatID(ctx, chatID, 10, &cursor)
	if err != nil || len(nextPage) != 2 || nextPage[0].ID != messages[1].ID || nextPage[1].ID != messages[0].ID {
		t.Fatalf("next page = %+v, error=%v", nextPage, err)
	}
	badCursor := uuid.New()
	if _, err := repo.GetByChatID(ctx, chatID, 10, &badCursor); err == nil {
		t.Fatal("GetByChatID() bad cursor error = nil")
	}

	unread, err := repo.CountUnreadForChat(ctx, chatID, currentUser, messages[0].ID)
	if err != nil || unread != 2 {
		t.Fatalf("CountUnreadForChat() = %d, %v", unread, err)
	}
	if _, err := repo.CountUnreadForChat(ctx, chatID, currentUser, badCursor); err == nil {
		t.Fatal("CountUnreadForChat() bad cursor error = nil")
	}

	if err := repo.Delete(ctx, generated.ID); err != nil {
		t.Fatal(err)
	}
	if generated.DeletedAt == nil {
		t.Fatal("Delete() did not set DeletedAt")
	}
	latestID, err = repo.GetLastMessageIDForChat(ctx, chatID)
	if err != nil || latestID != messages[2].ID {
		t.Fatalf("last visible message = %s, %v", latestID, err)
	}
	if err := repo.Delete(ctx, uuid.New()); err == nil {
		t.Fatal("Delete() missing message error = nil")
	}
}

func TestMessageRepositoryEmptyChatHasNoLastMessage(t *testing.T) {
	repo := NewMessageRepository()
	id, err := repo.GetLastMessageIDForChat(context.Background(), uuid.New())
	if err != nil || id != uuid.Nil {
		t.Fatalf("GetLastMessageIDForChat() = %s, %v", id, err)
	}
	if _, err := repo.GetByID(context.Background(), uuid.New()); err == nil {
		t.Fatal("GetByID() missing message error = nil")
	}
}
