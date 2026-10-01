package service

import (
	"context"
	"testing"

	"messenger/internal/model"
	"messenger/storage/memory"

	"github.com/google/uuid"
)

func TestDeleteMessageAllowsOnlySenderAndSoftDeletes(t *testing.T) {
	ctx := context.Background()
	messageRepo := memory.NewMessageRepository()
	membersRepo := memory.NeMembersRepository()
	chatID, senderID, otherID := uuid.New(), uuid.New(), uuid.New()
	message := &model.Message{ID: uuid.New(), ChatID: chatID, SenderID: senderID}

	if _, err := messageRepo.Create(ctx, message); err != nil {
		t.Fatal(err)
	}
	for _, userID := range []uuid.UUID{senderID, otherID} {
		if err := membersRepo.AddMember(ctx, &model.ChatMember{ChatID: chatID, UserID: userID}); err != nil {
			t.Fatal(err)
		}
	}

	service := &ChatService{messageRepo: messageRepo, membersRepo: membersRepo}
	if err := service.DeleteMessage(ctx, otherID, message.ID); err == nil {
		t.Fatal("DeleteMessage() by non-sender error = nil")
	}
	if message.DeletedAt != nil {
		t.Fatal("message was deleted by non-sender")
	}
	if err := service.DeleteMessage(ctx, senderID, message.ID); err != nil {
		t.Fatalf("DeleteMessage() by sender error = %v", err)
	}
	if message.DeletedAt == nil {
		t.Fatal("message DeletedAt was not set")
	}
}

func TestDeleteMessageRequiresChatMembership(t *testing.T) {
	ctx := context.Background()
	messageRepo := memory.NewMessageRepository()
	membersRepo := memory.NeMembersRepository()
	senderID := uuid.New()
	message := &model.Message{ID: uuid.New(), ChatID: uuid.New(), SenderID: senderID}
	if _, err := messageRepo.Create(ctx, message); err != nil {
		t.Fatal(err)
	}

	service := &ChatService{messageRepo: messageRepo, membersRepo: membersRepo}
	if err := service.DeleteMessage(ctx, senderID, message.ID); err == nil {
		t.Fatal("DeleteMessage() without membership error = nil")
	}
	if message.DeletedAt != nil {
		t.Fatal("message was deleted without chat membership")
	}
}
