package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"messenger/internal/model"
	storagegorm "messenger/storage/gorm"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type notificationCall struct {
	receiverID uuid.UUID
	senderID   uuid.UUID
	chatID     string
}

type notificationSpy struct{ calls chan notificationCall }

func (s *notificationSpy) SendNewMessage(_ context.Context, receiverID, senderID uuid.UUID, _, _, chatID string) error {
	s.calls <- notificationCall{receiverID: receiverID, senderID: senderID, chatID: chatID}
	return nil
}

func (*notificationSpy) Close() error { return nil }

func TestSendMessageNotifiesOnlyTheOtherPrivateChatMember(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.NewString())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		`CREATE TABLE chats (id text PRIMARY KEY, is_group numeric, title text, photo_url text, private_key text, created_at datetime)`,
		`CREATE TABLE chat_members (chat_id text, user_id text, role text, joined_at datetime, left_at datetime, muted numeric, last_read_message_id text, PRIMARY KEY (chat_id, user_id))`,
		`CREATE TABLE messages (id text PRIMARY KEY, chat_id text, sender_id text, text text, reply_to_id text, created_at datetime, edited_at datetime, deleted_at datetime, read_at datetime)`,
		`CREATE TABLE attachments (id text PRIMARY KEY, message_id text, url text, type text, file_name text, file_size integer)`,
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}

	base := storagegorm.NewRepository(db)
	chats := storagegorm.NewChatRepository(base)
	members := storagegorm.NewChatMemberRepository(base)
	messages := storagegorm.NewMessageRepository(base)
	attachments := storagegorm.NewAttachmentRepository(base)
	chatID, senderID, receiverID := uuid.New(), uuid.New(), uuid.New()
	if err := chats.Create(context.Background(), &model.Chat{ID: chatID}); err != nil {
		t.Fatal(err)
	}
	for _, userID := range []uuid.UUID{senderID, receiverID} {
		if err := members.AddMember(context.Background(), &model.ChatMember{ChatID: chatID, UserID: userID, Role: model.Member}); err != nil {
			t.Fatal(err)
		}
	}

	spy := &notificationSpy{calls: make(chan notificationCall, 2)}
	svc := &ChatService{chatRepo: chats, membersRepo: members, messageRepo: messages, attachmentRepo: attachments, notClient: spy}
	text := "hello"
	if _, err := svc.SendMessage(context.Background(), senderID, model.SendMessageInput{ChatID: chatID, SenderID: senderID, Text: &text}); err != nil {
		t.Fatal(err)
	}

	select {
	case call := <-spy.calls:
		if call.receiverID != receiverID || call.senderID != senderID || call.chatID != chatID.String() {
			t.Fatalf("notification routed incorrectly: %+v", call)
		}
	case <-time.After(time.Second):
		t.Fatal("notification was not delivered")
	}
	select {
	case extra := <-spy.calls:
		t.Fatalf("unexpected extra notification: %+v", extra)
	case <-time.After(50 * time.Millisecond):
	}
}
