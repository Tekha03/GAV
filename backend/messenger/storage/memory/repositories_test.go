package memory

import (
	"context"
	"errors"
	"messenger/internal/model"
	"messenger/internal/repository"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestChatRepository(t *testing.T) {
	repo := NewChatRepository()
	ctx := context.Background()
	chat := &model.Chat{ID: uuid.New(), Title: "dogs"}

	if err := repo.Create(ctx, chat); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := repo.Create(ctx, chat); !errors.Is(err, repository.ErrChatExists) {
		t.Fatalf("expected ErrChatExists, got %v", err)
	}

	stored, err := repo.GetByID(ctx, chat.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if stored.Title != "dogs" {
		t.Fatalf("unexpected title: %s", stored.Title)
	}

	if err := repo.UpdateTitle(ctx, chat.ID, "new"); err != nil {
		t.Fatalf("UpdateTitle returned error: %v", err)
	}
	if err := repo.UpdatePhoto(ctx, chat.ID, "avatar.png"); err != nil {
		t.Fatalf("UpdatePhoto returned error: %v", err)
	}
	stored, _ = repo.GetByID(ctx, chat.ID)
	if stored.Title != "new" || stored.PhotoURL != "avatar.png" {
		t.Fatalf("unexpected chat state: %+v", stored)
	}

	if err := repo.Delete(ctx, chat.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := repo.GetByID(ctx, chat.ID); !errors.Is(err, repository.ErrChatNotFound) {
		t.Fatalf("expected ErrChatNotFound, got %v", err)
	}
}

func TestMembersRepository(t *testing.T) {
	repo := NewMembersRepository()
	ctx := context.Background()
	chatID := uuid.New()
	userID := uuid.New()
	member := &model.ChatMember{ChatID: chatID, UserID: userID, Role: model.Member, LastReadMessageID: uuid.New()}

	if err := repo.AddMember(ctx, member); err != nil {
		t.Fatalf("AddMember returned error: %v", err)
	}
	if err := repo.AddMember(ctx, member); !errors.Is(err, repository.ErrMemberExists) {
		t.Fatalf("expected ErrMemberExists, got %v", err)
	}

	members, err := repo.GetMembers(ctx, chatID)
	if err != nil {
		t.Fatalf("GetMembers returned error: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("unexpected members count: %d", len(members))
	}

	admin := model.Admin
	if err := repo.UpdateRole(ctx, chatID, userID, &admin); err != nil {
		t.Fatalf("UpdateRole returned error: %v", err)
	}
	if err := repo.SetMuted(ctx, chatID, userID, true); err != nil {
		t.Fatalf("SetMuted returned error: %v", err)
	}
	members, _ = repo.GetMembers(ctx, chatID)
	if members[0].Role != model.Admin || !members[0].Muted {
		t.Fatalf("unexpected member state: %+v", members[0])
	}

	lastRead, err := repo.GetLastReadMessageID(ctx, chatID, userID)
	if err != nil {
		t.Fatalf("GetLastReadMessageID returned error: %v", err)
	}
	if lastRead != member.LastReadMessageID {
		t.Fatalf("unexpected lastRead: %s", lastRead)
	}

	chatIDs, err := repo.GetUserChats(ctx, userID)
	if err != nil {
		t.Fatalf("GetUserChats returned error: %v", err)
	}
	if len(chatIDs) != 1 || chatIDs[0] != chatID {
		t.Fatalf("unexpected chatIDs: %+v", chatIDs)
	}

	if err := repo.RemoveMember(ctx, userID, chatID); err != nil {
		t.Fatalf("RemoveMember returned error: %v", err)
	}
	if err := repo.RemoveMember(ctx, userID, chatID); !errors.Is(err, repository.ErrMemberNotFound) {
		t.Fatalf("expected ErrMemberNotFound, got %v", err)
	}
}

func TestMessageRepository(t *testing.T) {
	repo := NewMessageRepository()
	ctx := context.Background()
	chatID := uuid.New()
	userID := uuid.New()
	text := "first"
	first := &model.Message{ID: uuid.New(), ChatID: chatID, SenderID: userID, Text: &text}
	secondText := "second"
	second := &model.Message{ChatID: chatID, SenderID: userID, Text: &secondText}

	if _, err := repo.Create(ctx, first); err != nil {
		t.Fatalf("Create first returned error: %v", err)
	}
	if _, err := repo.Create(ctx, first); !errors.Is(err, repository.ErrMessageExists) {
		t.Fatalf("expected ErrMessageExists, got %v", err)
	}
	if _, err := repo.Create(ctx, second); err != nil {
		t.Fatalf("Create second returned error: %v", err)
	}

	updated := "updated"
	if err := repo.UpdateText(ctx, first.ID, updated); err != nil {
		t.Fatalf("UpdateText returned error: %v", err)
	}
	stored, err := repo.GetByID(ctx, first.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if stored.Text == nil || *stored.Text != updated || stored.EditedAt == nil {
		t.Fatalf("unexpected stored message: %+v", stored)
	}

	second.CreatedAt = time.Now().Add(2 * time.Minute)
	messages, err := repo.GetByChatID(ctx, chatID, 10, nil)
	if err != nil {
		t.Fatalf("GetByChatID returned error: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("unexpected messages count: %d", len(messages))
	}
	if messages[0].ID == messages[1].ID {
		t.Fatal("messages should be distinct")
	}

	cursor := first.ID
	filtered, err := repo.GetByChatID(ctx, chatID, 10, &cursor)
	if err != nil {
		t.Fatalf("GetByChatID with cursor returned error: %v", err)
	}
	if len(filtered) > len(messages) {
		t.Fatalf("cursor should not increase result size: %d > %d", len(filtered), len(messages))
	}

	limited, err := repo.GetByChatID(ctx, chatID, 1, nil)
	if err != nil {
		t.Fatalf("GetByChatID with limit returned error: %v", err)
	}
	if len(limited) != 1 {
		t.Fatalf("unexpected limited messages count: %d", len(limited))
	}

	if err := repo.Delete(ctx, second.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	messages, err = repo.GetByChatID(ctx, chatID, 10, nil)
	if err != nil {
		t.Fatalf("GetByChatID returned error: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("deleted message should be hidden, got %d", len(messages))
	}
	if err := repo.UpdateReadAtForChat(ctx, chatID, userID, time.Now()); err != nil {
		t.Fatalf("UpdateReadAtForChat returned error: %v", err)
	}
}

func TestAttachmentRepository(t *testing.T) {
	repo := NewAttachmentRepository()
	ctx := context.Background()
	messageID := uuid.New()
	attachment := &model.Attachment{ID: uuid.New(), MessageID: messageID, URL: "file", Type: model.AttachmentFile}

	if err := repo.Create(ctx, attachment); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := repo.Create(ctx, attachment); !errors.Is(err, ErrAttachmentExist) {
		t.Fatalf("expected ErrAttachmentExist, got %v", err)
	}
	stored, err := repo.GetByID(ctx, attachment.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if stored.MessageID != messageID {
		t.Fatalf("unexpected attachment: %+v", stored)
	}

	attachments, err := repo.GetByMessage(ctx, messageID)
	if err != nil {
		t.Fatalf("GetByMessage returned error: %v", err)
	}
	if len(attachments) != 1 {
		t.Fatalf("unexpected attachments count: %d", len(attachments))
	}

	if err := repo.DeleteByMessage(ctx, messageID); err != nil {
		t.Fatalf("DeleteByMessage returned error: %v", err)
	}
	attachments, err = repo.GetByMessage(ctx, messageID)
	if err != nil {
		t.Fatalf("GetByMessage returned error: %v", err)
	}
	if len(attachments) != 0 {
		t.Fatalf("attachments should be deleted, got %d", len(attachments))
	}
	if err := repo.Delete(ctx, attachment.ID); !errors.Is(err, ErrAttachmentNotFound) {
		t.Fatalf("expected ErrAttachmentNotFound, got %v", err)
	}
}

func TestPinnedReactionAndTypingRepositories(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()
	messageID := uuid.New()
	userID := uuid.New()

	pinnedRepo := NewPinnedRepo()
	if err := pinnedRepo.Pin(ctx, chatID, messageID); err != nil {
		t.Fatalf("Pin returned error: %v", err)
	}
	ids := pinnedRepo.GetByChatID(ctx, chatID)
	if len(ids) != 1 || ids[0] != messageID {
		t.Fatalf("unexpected pinned ids: %+v", ids)
	}
	if err := pinnedRepo.Unpin(ctx, chatID, messageID); err != nil {
		t.Fatalf("Unpin returned error: %v", err)
	}
	ids = pinnedRepo.GetByChatID(ctx, chatID)
	if len(ids) != 0 {
		t.Fatalf("unexpected pinned ids after unpin: %+v", ids)
	}

	reactionRepo := NewReactionRepository()
	reaction := &model.Reaction{MessageID: messageID, UserID: userID, Emoji: "🔥"}
	if err := reactionRepo.Add(ctx, reaction); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if err := reactionRepo.Add(ctx, reaction); !errors.Is(err, repository.ErrReactionExists) {
		t.Fatalf("expected ErrReactionExists, got %v", err)
	}
	if err := reactionRepo.Remove(ctx, messageID, userID); err != nil {
		t.Fatalf("Remove returned error: %v", err)
	}
	if err := reactionRepo.Remove(ctx, messageID, userID); !errors.Is(err, repository.ErrReactionNotFound) {
		t.Fatalf("expected ErrReactionNotFound, got %v", err)
	}

	typingRepo := NewTypingRepository()
	typingRepo.SetTyping(ctx, chatID, userID)
	users := typingRepo.GetTypingUsers(ctx, chatID, time.Minute)
	if len(users) != 1 || users[0] != userID {
		t.Fatalf("unexpected typing users: %+v", users)
	}
	typingRepo.Cleanup(ctx, -time.Second)
	users = typingRepo.GetTypingUsers(ctx, chatID, time.Minute)
	if len(users) != 0 {
		t.Fatalf("cleanup should remove user, got %+v", users)
	}
}
