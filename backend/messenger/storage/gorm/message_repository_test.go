package gorm

import (
	"context"
	"testing"
	"time"

	"messenger/internal/model"
	apperrors "shared/app_errors"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	gormpkg "gorm.io/gorm"
)

func TestMessageRepositoryGetByChatIDFirstPageUsesStableOrder(t *testing.T) {
	ctx := context.Background()
	repo := newTestMessageRepository(t)

	chatID := uuid.New()
	olderID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	lowerTieID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	higherTieID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	sameTime := time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC)

	insertTestMessage(t, repo, model.Message{ID: olderID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime.Add(-time.Minute)})
	insertTestMessage(t, repo, model.Message{ID: lowerTieID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime})
	insertTestMessage(t, repo, model.Message{ID: higherTieID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime})

	messages, err := repo.GetByChatID(ctx, chatID, 10, nil)
	if err != nil {
		t.Fatalf("GetByChatID() error = %v", err)
	}

	assertMessageIDs(t, messages, []uuid.UUID{higherTieID, lowerTieID, olderID})
}

func TestMessageRepositoryGetByChatIDBeforeCursorReturnsOlderMessages(t *testing.T) {
	ctx := context.Background()
	repo := newTestMessageRepository(t)

	chatID := uuid.New()
	newestID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	cursorID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	olderID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	baseTime := time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC)

	insertTestMessage(t, repo, model.Message{ID: olderID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: baseTime.Add(-2 * time.Minute)})
	insertTestMessage(t, repo, model.Message{ID: cursorID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: baseTime.Add(-time.Minute)})
	insertTestMessage(t, repo, model.Message{ID: newestID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: baseTime})

	messages, err := repo.GetByChatID(ctx, chatID, 10, &cursorID)
	if err != nil {
		t.Fatalf("GetByChatID() error = %v", err)
	}

	assertMessageIDs(t, messages, []uuid.UUID{olderID})
}

func TestMessageRepositoryGetByChatIDBeforeCursorHandlesTimestampTies(t *testing.T) {
	ctx := context.Background()
	repo := newTestMessageRepository(t)

	chatID := uuid.New()
	newerTieID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	cursorID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	olderTieID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	olderByTimeID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	sameTime := time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC)

	insertTestMessage(t, repo, model.Message{ID: newerTieID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime})
	insertTestMessage(t, repo, model.Message{ID: cursorID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime})
	insertTestMessage(t, repo, model.Message{ID: olderTieID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime})
	insertTestMessage(t, repo, model.Message{ID: olderByTimeID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime.Add(-time.Minute)})

	messages, err := repo.GetByChatID(ctx, chatID, 10, &cursorID)
	if err != nil {
		t.Fatalf("GetByChatID() error = %v", err)
	}

	assertMessageIDs(t, messages, []uuid.UUID{olderTieID, olderByTimeID})
}

func TestMessageRepositoryGetByChatIDRejectsCursorFromAnotherChat(t *testing.T) {
	ctx := context.Background()
	repo := newTestMessageRepository(t)

	chatID := uuid.New()
	otherChatID := uuid.New()
	cursorID := uuid.New()

	insertTestMessage(t, repo, model.Message{
		ID:        cursorID,
		ChatID:    otherChatID,
		SenderID:  uuid.New(),
		CreatedAt: time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC),
	})

	_, err := repo.GetByChatID(ctx, chatID, 10, &cursorID)
	if err == nil {
		t.Fatal("GetByChatID() error = nil, want cursor message not found")
	}

	appErr := apperrors.Normalize(err)
	if appErr.Code != apperrors.MessageNotFound.Code {
		t.Fatalf("GetByChatID() error code = %s, want %s", appErr.Code, apperrors.MessageNotFound.Code)
	}
}

func TestMessageRepositoryUpdateLastReadMessageForChatStoresLatestMessageID(t *testing.T) {
	ctx := context.Background()
	repo := newTestMessageRepository(t)

	chatID := uuid.New()
	userID := uuid.New()
	olderID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	newerID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	baseTime := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

	insertTestChatMember(t, repo, chatID, userID)
	insertTestMessage(t, repo, model.Message{ID: olderID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: baseTime})
	insertTestMessage(t, repo, model.Message{ID: newerID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: baseTime.Add(time.Minute)})

	if err := repo.UpdateLastReadMessageForChat(ctx, chatID, userID); err != nil {
		t.Fatalf("UpdateLastReadMessageForChat() error = %v", err)
	}

	got := getLastReadMessageID(t, repo, chatID, userID)
	if got != newerID {
		t.Fatalf("last_read_message_id = %s, want %s", got, newerID)
	}
}

func TestMessageRepositoryUpdateLastReadMessageForChatUsesStableOrderForTimestampTie(t *testing.T) {
	ctx := context.Background()
	repo := newTestMessageRepository(t)

	chatID := uuid.New()
	userID := uuid.New()
	lowerTieID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	higherTieID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	sameTime := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

	insertTestChatMember(t, repo, chatID, userID)
	insertTestMessage(t, repo, model.Message{ID: lowerTieID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime})
	insertTestMessage(t, repo, model.Message{ID: higherTieID, ChatID: chatID, SenderID: uuid.New(), CreatedAt: sameTime})

	if err := repo.UpdateLastReadMessageForChat(ctx, chatID, userID); err != nil {
		t.Fatalf("UpdateLastReadMessageForChat() error = %v", err)
	}

	got := getLastReadMessageID(t, repo, chatID, userID)
	if got != higherTieID {
		t.Fatalf("last_read_message_id = %s, want %s", got, higherTieID)
	}
}

func TestMessageRepositoryUpdateLastReadMessageForChatEmptyChatKeepsZeroMessageID(t *testing.T) {
	ctx := context.Background()
	repo := newTestMessageRepository(t)

	chatID := uuid.New()
	userID := uuid.New()
	insertTestChatMember(t, repo, chatID, userID)

	if err := repo.UpdateLastReadMessageForChat(ctx, chatID, userID); err != nil {
		t.Fatalf("UpdateLastReadMessageForChat() error = %v", err)
	}

	got := getLastReadMessageID(t, repo, chatID, userID)
	if got != uuid.Nil {
		t.Fatalf("last_read_message_id = %s, want nil UUID", got)
	}
}

func TestMessageRepositoryUpdateLastReadMessageForChatRejectsNonMember(t *testing.T) {
	ctx := context.Background()
	repo := newTestMessageRepository(t)

	err := repo.UpdateLastReadMessageForChat(ctx, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("UpdateLastReadMessageForChat() error = nil, want chat member not found")
	}

	appErr := apperrors.Normalize(err)
	if appErr.Code != apperrors.ChatMemberNotFound.Code {
		t.Fatalf("UpdateLastReadMessageForChat() error code = %s, want %s", appErr.Code, apperrors.ChatMemberNotFound.Code)
	}
}

func newTestMessageRepository(t *testing.T) *MessageRepository {
	t.Helper()

	db, err := gormpkg.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gormpkg.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	if err := db.Exec(`CREATE TABLE messages (
		id text PRIMARY KEY,
		chat_id text NOT NULL,
		sender_id text NOT NULL,
		text text,
		reply_to_id text,
		created_at datetime,
		edited_at datetime,
		deleted_at datetime,
		read_at datetime
	)`).Error; err != nil {
		t.Fatalf("create messages table: %v", err)
	}
	if err := db.Exec(`CREATE TABLE chat_members (
		chat_id text NOT NULL,
		user_id text NOT NULL,
		role text NOT NULL,
		joined_at datetime,
		left_at datetime,
		muted numeric NOT NULL DEFAULT false,
		last_read_message_id text,
		PRIMARY KEY (chat_id, user_id)
	)`).Error; err != nil {
		t.Fatalf("create chat_members table: %v", err)
	}

	return NewMessageRepository(NewRepository(db)).(*MessageRepository)
}

func insertTestMessage(t *testing.T, repo *MessageRepository, message model.Message) {
	t.Helper()

	if err := repo.repo.db.Create(&message).Error; err != nil {
		t.Fatalf("insert message %s: %v", message.ID, err)
	}
}

func insertTestChatMember(t *testing.T, repo *MessageRepository, chatID, userID uuid.UUID) {
	t.Helper()

	member := model.ChatMember{
		ChatID: chatID,
		UserID: userID,
		Role:   model.Member,
	}
	if err := repo.repo.db.Create(&member).Error; err != nil {
		t.Fatalf("insert chat member %s/%s: %v", chatID, userID, err)
	}
}

func getLastReadMessageID(t *testing.T, repo *MessageRepository, chatID, userID uuid.UUID) uuid.UUID {
	t.Helper()

	var member model.ChatMember
	if err := repo.repo.db.
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		First(&member).Error; err != nil {
		t.Fatalf("get chat member %s/%s: %v", chatID, userID, err)
	}
	return member.LastReadMessageID
}

func assertMessageIDs(t *testing.T, messages []*model.Message, want []uuid.UUID) {
	t.Helper()

	if len(messages) != len(want) {
		t.Fatalf("messages len = %d, want %d", len(messages), len(want))
	}

	for i, message := range messages {
		if message.ID != want[i] {
			t.Fatalf("messages[%d].ID = %s, want %s", i, message.ID, want[i])
		}
	}
}
