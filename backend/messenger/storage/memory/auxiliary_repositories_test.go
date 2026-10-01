package memory

import (
	"context"
	"testing"
	"time"

	"messenger/internal/model"

	"github.com/google/uuid"
)

func TestChatRepositoryLifecycleAndPrivateKeyUniqueness(t *testing.T) {
	ctx := context.Background()
	repo := NewChatRepository()
	privateKey := "first:second"
	chat := &model.Chat{PrivateKey: &privateKey, Title: "Old"}
	if err := repo.Create(ctx, chat); err != nil || chat.ID == uuid.Nil {
		t.Fatalf("Create() id=%s error=%v", chat.ID, err)
	}
	if err := repo.Create(ctx, &model.Chat{PrivateKey: &privateKey}); err == nil {
		t.Fatal("duplicate private key error = nil")
	}
	if err := repo.Create(ctx, chat); err == nil {
		t.Fatal("duplicate id error = nil")
	}
	if err := repo.UpdateTitle(ctx, chat.ID, "New"); err != nil {
		t.Fatal(err)
	}
	if err := repo.UpdatePhoto(ctx, chat.ID, "photo.jpg"); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.GetByID(ctx, chat.ID)
	if err != nil || loaded.Title != "New" || loaded.PhotoURL != "photo.jpg" {
		t.Fatalf("GetByID() = %+v, %v", loaded, err)
	}
	byKey, err := repo.GetByPrivateKey(ctx, privateKey)
	if err != nil || byKey != chat {
		t.Fatalf("GetByPrivateKey() = %+v, %v", byKey, err)
	}
	if missing, err := repo.GetByPrivateKey(ctx, "missing"); err != nil || missing != nil {
		t.Fatalf("GetByPrivateKey(missing) = %+v, %v", missing, err)
	}
	if err := repo.Delete(ctx, chat.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(ctx, chat.ID); err == nil {
		t.Fatal("GetByID() deleted chat error = nil")
	}
	if err := repo.Delete(ctx, chat.ID); err == nil {
		t.Fatal("Delete() missing chat error = nil")
	}
	if err := repo.UpdateTitle(ctx, chat.ID, "x"); err == nil {
		t.Fatal("UpdateTitle() missing chat error = nil")
	}
	if err := repo.UpdatePhoto(ctx, chat.ID, "x"); err == nil {
		t.Fatal("UpdatePhoto() missing chat error = nil")
	}
}

func TestAttachmentRepositoryLifecycleAndBatch(t *testing.T) {
	ctx := context.Background()
	repo := NewAttachmentRepository()
	firstMessage, secondMessage := uuid.New(), uuid.New()
	first := &model.Attachment{MessageID: firstMessage, URL: "first.jpg"}
	if err := repo.Create(ctx, first); err != nil || first.ID == uuid.Nil {
		t.Fatalf("Create() = %+v, %v", first, err)
	}
	if err := repo.Create(ctx, first); err == nil {
		t.Fatal("duplicate Create() error = nil")
	}
	batch := []model.Attachment{
		{MessageID: firstMessage, URL: "second.jpg"},
		{MessageID: secondMessage, URL: "third.jpg"},
	}
	if err := repo.CreateBatch(ctx, batch); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.GetByID(ctx, first.ID)
	if err != nil || loaded.URL != first.URL {
		t.Fatalf("GetByID() = %+v, %v", loaded, err)
	}
	forMessage, err := repo.GetByMessage(ctx, firstMessage)
	if err != nil || len(forMessage) != 2 {
		t.Fatalf("GetByMessage() = %+v, %v", forMessage, err)
	}
	if err := repo.Delete(ctx, first.ID); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, first.ID); err == nil {
		t.Fatal("Delete() missing attachment error = nil")
	}
	if _, err := repo.GetByID(ctx, first.ID); err == nil {
		t.Fatal("GetByID() missing attachment error = nil")
	}
	if err := repo.DeleteByMessage(ctx, firstMessage); err != nil {
		t.Fatal(err)
	}
	if remaining, _ := repo.GetByMessage(ctx, firstMessage); len(remaining) != 0 {
		t.Fatalf("attachments remained = %+v", remaining)
	}
}

func TestReactionPinnedAndTypingRepositories(t *testing.T) {
	ctx := context.Background()
	messageID, chatID, userID := uuid.New(), uuid.New(), uuid.New()

	reactions := NewReactionRepository()
	reaction := &model.Reaction{MessageID: messageID, UserID: userID, Emoji: "👍"}
	if err := reactions.Add(ctx, reaction); err != nil || reaction.ID == uuid.Nil {
		t.Fatalf("Add() = %+v, %v", reaction, err)
	}
	if err := reactions.Add(ctx, reaction); err == nil {
		t.Fatal("duplicate reaction error = nil")
	}
	if err := reactions.Remove(ctx, messageID, userID); err != nil {
		t.Fatal(err)
	}
	if err := reactions.Remove(ctx, messageID, userID); err == nil {
		t.Fatal("missing reaction error = nil")
	}
	if err := reactions.Remove(ctx, uuid.New(), userID); err == nil {
		t.Fatal("missing message reactions error = nil")
	}

	pinned := NewPinnedRepo()
	if err := pinned.Pin(ctx, chatID, messageID); err != nil {
		t.Fatal(err)
	}
	if ids, err := pinned.GetByChatID(ctx, chatID); err != nil || len(ids) != 1 || ids[0] != messageID {
		t.Fatalf("GetByChatID() = %+v, %v", ids, err)
	}
	if err := pinned.Unpin(ctx, chatID, messageID); err != nil {
		t.Fatal(err)
	}
	if err := pinned.Unpin(ctx, uuid.New(), messageID); err != nil {
		t.Fatal(err)
	}

	typing := NewTypingRepository()
	if err := typing.SetTyping(ctx, chatID, userID); err != nil {
		t.Fatal(err)
	}
	if users, err := typing.GetTypingUsers(ctx, chatID, time.Minute); err != nil || len(users) != 1 || users[0] != userID {
		t.Fatalf("GetTypingUsers() = %+v, %v", users, err)
	}
	if users, err := typing.GetTypingUsers(ctx, uuid.New(), time.Minute); err != nil || len(users) != 0 {
		t.Fatalf("GetTypingUsers(missing) = %+v, %v", users, err)
	}
	if err := typing.Cleanup(ctx, -time.Second); err != nil {
		t.Fatal(err)
	}
	if users, _ := typing.GetTypingUsers(ctx, chatID, time.Minute); len(users) != 0 {
		t.Fatalf("typing users after cleanup = %+v", users)
	}
}
