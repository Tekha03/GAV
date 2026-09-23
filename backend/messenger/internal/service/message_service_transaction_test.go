package service

import (
	"context"
	"fmt"
	"messenger/internal/model"
	storagegorm "messenger/storage/gorm"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSendMessageRollsBackWhenOutboxWriteFails(t *testing.T) {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	schema := []string{
		`CREATE TABLE chats (
			id text PRIMARY KEY,
			is_group numeric NOT NULL DEFAULT false,
			title text,
			photo_url text,
			private_key text,
			created_at datetime
		)`,
		`CREATE TABLE chat_members (
			chat_id text NOT NULL,
			user_id text NOT NULL,
			role text NOT NULL,
			joined_at datetime,
			left_at datetime,
			muted numeric NOT NULL DEFAULT false,
			last_read_message_id text,
			PRIMARY KEY (chat_id, user_id)
		)`,
		`CREATE TABLE messages (
			id text PRIMARY KEY,
			chat_id text NOT NULL,
			sender_id text NOT NULL,
			text text,
			reply_to_id text,
			created_at datetime,
			edited_at datetime,
			deleted_at datetime,
			read_at datetime
		)`,
		`CREATE TABLE attachments (
			id text PRIMARY KEY,
			message_id text NOT NULL,
			url text NOT NULL,
			type text NOT NULL,
			file_name text,
			file_size integer NOT NULL DEFAULT 0
		)`,
	}
	for _, statement := range schema {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("create test schema: %v", err)
		}
	}

	baseRepo := storagegorm.NewRepository(db)
	chatRepo := storagegorm.NewChatRepository(baseRepo)
	membersRepo := storagegorm.NewChatMemberRepository(baseRepo)
	messageRepo := storagegorm.NewMessageRepository(baseRepo)
	attachmentRepo := storagegorm.NewAttachmentRepository(baseRepo)
	outboxRepo := storagegorm.NewOutboxRepository(baseRepo)

	chatID := uuid.New()
	senderID := uuid.New()
	receiverID := uuid.New()
	if err := chatRepo.Create(context.Background(), &model.Chat{ID: chatID}); err != nil {
		t.Fatalf("create chat: %v", err)
	}
	for _, userID := range []uuid.UUID{senderID, receiverID} {
		if err := membersRepo.AddMember(context.Background(), &model.ChatMember{
			ChatID: chatID,
			UserID: userID,
			Role:   model.Member,
		}); err != nil {
			t.Fatalf("add member: %v", err)
		}
	}

	svc := &ChatService{
		transactionManager: baseRepo,
		chatRepo:           chatRepo,
		membersRepo:        membersRepo,
		messageRepo:        messageRepo,
		attachmentRepo:     attachmentRepo,
		outboxRepo:         outboxRepo,
	}
	text := "hello"
	_, err = svc.SendMessage(context.Background(), senderID, model.SendMessageInput{
		ChatID:   chatID,
		SenderID: senderID,
		Text:     &text,
		Attachments: []model.AttachmentInput{{
			Type:     model.AttachmentImage,
			URL:      "https://example.test/image.jpg",
			FileName: "image.jpg",
			FileSize: 10,
		}},
	})
	if err == nil {
		t.Fatal("SendMessage() error = nil, want outbox write error")
	}

	for table, destination := range map[string]any{
		"messages":    &[]model.Message{},
		"attachments": &[]model.Attachment{},
	} {
		var count int64
		if err := db.Model(destination).Count(&count).Error; err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("%s count = %d, want 0 after rollback", table, count)
		}
	}
}
