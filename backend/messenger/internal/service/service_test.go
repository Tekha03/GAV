package service

import (
	"context"
	"encoding/json"
	"errors"
	"messenger/internal/constatnts"
	appErrors "messenger/internal/errors"
	"messenger/internal/model"
	"messenger/internal/repository"
	memory "messenger/storage/memory"
	"shared/events"
	"testing"
	"time"

	"github.com/google/uuid"
)

type stubProducer struct {
	events []events.Event
	err    error
}

func (p *stubProducer) PublishEvent(event events.Event) error {
	if p.err != nil {
		return p.err
	}
	p.events = append(p.events, event)
	return nil
}

type nilMembersRepo struct {
	repository.ChatMemberRepository
}

func (r nilMembersRepo) GetMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error) {
	return nil, nil
}

type nilUserChatsRepo struct {
	repository.ChatMemberRepository
}

func (r nilUserChatsRepo) GetUserChats(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

type nilMessageRepo struct{ repository.MessageRepository }

func (r nilMessageRepo) GetByID(ctx context.Context, messageID uuid.UUID) (*model.Message, error) {
	return nil, nil
}

func (r nilMessageRepo) GetByChatID(ctx context.Context, chatID uuid.UUID, limit int, cursorID *uuid.UUID) ([]*model.Message, error) {
	return nil, nil
}

func newTestService() (*ChatService, *stubProducer) {
	producer := &stubProducer{}
	return &ChatService{
		chatRepo:       memory.NewChatRepository(),
		membersRepo:    memory.NewMembersRepository(),
		messageRepo:    memory.NewMessageRepository(),
		attachmentRepo: memory.NewAttachmentRepository(),
		reactionRepo:   memory.NewReactionRepository(),
		pinnedRepo:     memory.NewPinnedRepo(),
		typingRepo:     memory.NewTypingRepository(),
		producer:       producer,
	}, producer
}

func createChatWithMembers(t *testing.T, svc *ChatService, users ...uuid.UUID) uuid.UUID {
	t.Helper()

	chat := &model.Chat{ID: uuid.New(), CreatedAt: time.Now()}
	if err := svc.chatRepo.Create(context.Background(), chat); err != nil {
		t.Fatalf("create chat: %v", err)
	}

	for _, userID := range users {
		member := &model.ChatMember{ChatID: chat.ID, UserID: userID, Role: model.Member, JoinedAt: time.Now()}
		if err := svc.membersRepo.AddMember(context.Background(), member); err != nil {
			t.Fatalf("add member: %v", err)
		}
	}

	return chat.ID
}

func createMessage(t *testing.T, svc *ChatService, chatID, senderID uuid.UUID, text string) *model.Message {
	t.Helper()
	msg := &model.Message{ChatID: chatID, SenderID: senderID, Text: &text}
	id, err := svc.messageRepo.Create(context.Background(), msg)
	if err != nil {
		t.Fatalf("create message: %v", err)
	}
	msg.ID = id
	return msg
}

func TestCreatePrivateChat(t *testing.T) {
	svc, producer := newTestService()
	ctx := context.Background()
	firstUser := uuid.New()
	secondUser := uuid.New()

	chat, err := svc.CreatePrivateChat(ctx, firstUser, secondUser)
	if err != nil {
		t.Fatalf("CreatePrivateChat returned error: %v", err)
	}
	if chat == nil {
		t.Fatal("chat is nil")
	}
	if chat.IsGroup {
		t.Fatal("private chat marked as group")
	}

	members, err := svc.GetChatMembers(ctx, chat.ID)
	if err != nil {
		t.Fatalf("GetChatMembers returned error: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("unexpected members count: %d", len(members))
	}

	if len(producer.events) != 1 {
		t.Fatalf("unexpected events count: %d", len(producer.events))
	}
	if producer.events[0].EventType != events.EventTypeChatCreated {
		t.Fatalf("unexpected event type: %s", producer.events[0].EventType)
	}

	var payload events.ChatCreatedData
	if err := json.Unmarshal(producer.events[0].Data, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload.ChatID != chat.ID {
		t.Fatalf("unexpected chat id in payload: %s", payload.ChatID)
	}
	if len(payload.Members) != 2 {
		t.Fatalf("unexpected members in payload: %d", len(payload.Members))
	}

	_, err = svc.CreatePrivateChat(ctx, firstUser, firstUser)
	if !errors.Is(err, appErrors.ErrChatWithYourSelf) {
		t.Fatalf("expected ErrChatWithYourSelf, got %v", err)
	}
}

func TestCreateGroupChat(t *testing.T) {
	svc, producer := newTestService()
	ctx := context.Background()
	creatorID := uuid.New()
	memberID := uuid.New()

	chat, err := svc.CreateGroupChat(ctx, "dogs", creatorID, []uuid.UUID{creatorID, memberID, memberID})
	if err != nil {
		t.Fatalf("CreateGroupChat returned error: %v", err)
	}
	if chat == nil || !chat.IsGroup {
		t.Fatal("group chat was not created")
	}

	members, err := svc.GetChatMembers(ctx, chat.ID)
	if err != nil {
		t.Fatalf("GetChatMembers returned error: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("unexpected members count: %d", len(members))
	}

	roles := map[uuid.UUID]model.MemberRole{}
	for _, member := range members {
		roles[member.UserID] = member.Role
	}
	if roles[creatorID] != model.Admin {
		t.Fatalf("creator should be admin, got %s", roles[creatorID])
	}
	if roles[memberID] != model.Member {
		t.Fatalf("member should be regular member, got %s", roles[memberID])
	}

	if len(producer.events) != 1 || producer.events[0].EventType != events.EventTypeChatCreated {
		t.Fatalf("unexpected events: %+v", producer.events)
	}
}

func TestChatMembershipAndUpdates(t *testing.T) {
	svc, producer := newTestService()
	ctx := context.Background()
	ownerID := uuid.New()
	guestID := uuid.New()
	chatID := createChatWithMembers(t, svc, ownerID)

	if err := svc.AddMember(ctx, guestID, chatID); err != nil {
		t.Fatalf("AddMember returned error: %v", err)
	}
	members, err := svc.GetChatMembers(ctx, chatID)
	if err != nil {
		t.Fatalf("GetChatMembers returned error: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("unexpected members count after add: %d", len(members))
	}
	if len(producer.events) != 1 || producer.events[0].EventType != events.EventTypeChatMemberAdded {
		t.Fatalf("unexpected add member event: %+v", producer.events)
	}

	userChats, err := svc.GetUserChats(ctx, guestID)
	if err != nil {
		t.Fatalf("GetUserChats returned error: %v", err)
	}
	if len(userChats) != 1 || userChats[0].ID != chatID {
		t.Fatalf("unexpected user chats: %+v", userChats)
	}

	if err := svc.UpdateChatTitle(ctx, chatID, "new title"); err != nil {
		t.Fatalf("UpdateChatTitle returned error: %v", err)
	}
	if err := svc.UpdateChatPhoto(ctx, chatID, "avatar.png"); err != nil {
		t.Fatalf("UpdateChatPhoto returned error: %v", err)
	}
	chat, err := svc.GetChatByID(ctx, chatID)
	if err != nil {
		t.Fatalf("GetChatByID returned error: %v", err)
	}
	if chat.Title != "new title" || chat.PhotoURL != "avatar.png" {
		t.Fatalf("unexpected chat state: %+v", chat)
	}

	if err := svc.RemoveMember(ctx, guestID, chatID); err != nil {
		t.Fatalf("RemoveMember returned error: %v", err)
	}
	if len(producer.events) != 2 || producer.events[1].EventType != events.EventTypeChatMemberRemoved {
		t.Fatalf("unexpected remove member event: %+v", producer.events)
	}

	if err := svc.AddMember(ctx, guestID, chatID); err != nil {
		t.Fatalf("AddMember returned error: %v", err)
	}
	if err := svc.LeaveChat(ctx, guestID, chatID); err != nil {
		t.Fatalf("LeaveChat returned error: %v", err)
	}
	members, err = svc.GetChatMembers(ctx, chatID)
	if err != nil {
		t.Fatalf("GetChatMembers returned error: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("unexpected members count after leave: %d", len(members))
	}
}

func TestChatErrorBranches(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()

	_, err := svc.GetChatByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for missing chat")
	}

	nilMembersService, _ := newTestService()
	nilMembersService.membersRepo = nilMembersRepo{nilMembersService.membersRepo}
	_, err = nilMembersService.GetChatMembers(ctx, uuid.New())
	if !errors.Is(err, appErrors.ErrNoMembers) {
		t.Fatalf("expected ErrNoMembers, got %v", err)
	}

	nilChatsService, _ := newTestService()
	nilChatsService.membersRepo = nilUserChatsRepo{nilChatsService.membersRepo}
	_, err = nilChatsService.GetUserChats(ctx, uuid.New())
	if !errors.Is(err, appErrors.ErrNoChats) {
		t.Fatalf("expected ErrNoChats, got %v", err)
	}

	chatID := createChatWithMembers(t, svc, uuid.New())
	missingChatID := uuid.New()
	if err := svc.UpdateChatTitle(ctx, missingChatID, "x"); !errors.Is(err, appErrors.ErrTitleUpdate) {
		t.Fatalf("expected ErrTitleUpdate, got %v", err)
	}
	if err := svc.UpdateChatPhoto(ctx, missingChatID, "x"); !errors.Is(err, appErrors.ErrPhotoUpdate) {
		t.Fatalf("expected ErrPhotoUpdate, got %v", err)
	}
	if _, err := svc.GetChatByID(ctx, chatID); err != nil {
		t.Fatalf("existing chat should be available: %v", err)
	}
}

func TestSendMessageValidationAndSuccess(t *testing.T) {
	svc, producer := newTestService()
	ctx := context.Background()
	senderID := uuid.New()
	receiverID := uuid.New()
	chatID := createChatWithMembers(t, svc, senderID, receiverID)

	_, err := svc.SendMessage(ctx, model.SendMessageInput{ChatID: chatID, SenderID: senderID})
	if !errors.Is(err, appErrors.ErrEmptyMessage) {
		t.Fatalf("expected ErrEmptyMessage, got %v", err)
	}

	longText := make([]byte, constatnts.MaxMessageLength+1)
	for i := range longText {
		longText[i] = 'a'
	}
	tooLong := string(longText)
	_, err = svc.SendMessage(ctx, model.SendMessageInput{ChatID: chatID, SenderID: senderID, Text: &tooLong})
	if !errors.Is(err, appErrors.ErrTextOverLength) {
		t.Fatalf("expected ErrTextOverLength, got %v", err)
	}

	attachments := make([]model.AttachmentInput, constatnts.MaxAttachments+1)
	_, err = svc.SendMessage(ctx, model.SendMessageInput{ChatID: chatID, SenderID: senderID, Attachments: attachments})
	if !errors.Is(err, appErrors.ErrAttachmentsOverLength) {
		t.Fatalf("expected ErrAttachmentsOverLength, got %v", err)
	}

	otherChatID := createChatWithMembers(t, svc, senderID)
	replyText := "reply target"
	replyMessage := createMessage(t, svc, otherChatID, senderID, replyText)
	text := "hello"
	_, err = svc.SendMessage(ctx, model.SendMessageInput{ChatID: chatID, SenderID: senderID, Text: &text, ReplyToID: &replyMessage.ID})
	if !errors.Is(err, appErrors.ErrInvalidReply) {
		t.Fatalf("expected ErrInvalidReply, got %v", err)
	}

	body := "hello, dogs"
	message, err := svc.SendMessage(ctx, model.SendMessageInput{
		ChatID:   chatID,
		SenderID: senderID,
		Text:     &body,
		Attachments: []model.AttachmentInput{{
			Type:     model.AttachmentImage,
			URL:      "https://cdn/image.png",
			FileName: "image.png",
			FileSize: 42,
		}},
	})
	if err != nil {
		t.Fatalf("SendMessage returned error: %v", err)
	}
	if message == nil || message.ID == uuid.Nil {
		t.Fatal("message was not created")
	}

	storedAttachments, err := svc.attachmentRepo.GetByMessage(ctx, message.ID)
	if err != nil {
		t.Fatalf("GetByMessage returned error: %v", err)
	}
	if len(storedAttachments) != 1 {
		t.Fatalf("unexpected attachments count: %d", len(storedAttachments))
	}
	if len(producer.events) != 1 || producer.events[0].EventType != events.EventTypeMessageSent {
		t.Fatalf("unexpected producer events: %+v", producer.events)
	}

	var payload events.MessageSentData
	if err := json.Unmarshal(producer.events[0].Data, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload.MessageID != message.ID || payload.ChatID != chatID || payload.SenderID != senderID || payload.Text != body {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestSendMessageProducerErrorAndMissingMembers(t *testing.T) {
	svc, producer := newTestService()
	ctx := context.Background()
	senderID := uuid.New()
	chatID := createChatWithMembers(t, svc, senderID, uuid.New())
	text := "hello"
	producer.err = errors.New("kafka is down")

	_, err := svc.SendMessage(ctx, model.SendMessageInput{ChatID: chatID, SenderID: senderID, Text: &text})
	if err == nil || err.Error() != "kafka is down" {
		t.Fatalf("expected producer error, got %v", err)
	}

	nilMembersService, _ := newTestService()
	nilMembersService.chatRepo.Create(ctx, &model.Chat{ID: chatID})
	nilMembersService.membersRepo = nilMembersRepo{nilMembersService.membersRepo}
	_, err = nilMembersService.SendMessage(ctx, model.SendMessageInput{ChatID: chatID, SenderID: senderID, Text: &text})
	if !errors.Is(err, appErrors.ErrNoMembers) {
		t.Fatalf("expected ErrNoMembers, got %v", err)
	}
}

func TestEditDeleteGetMessagesAndMarkAsRead(t *testing.T) {
	svc, producer := newTestService()
	ctx := context.Background()
	senderID := uuid.New()
	receiverID := uuid.New()
	chatID := createChatWithMembers(t, svc, senderID, receiverID)

	text := "old"
	message := createMessage(t, svc, chatID, senderID, text)

	edited, err := svc.EditMessage(ctx, message.ID, "new")
	if err != nil {
		t.Fatalf("EditMessage returned error: %v", err)
	}
	if edited.Text == nil || *edited.Text != "new" {
		t.Fatalf("unexpected edited message: %+v", edited)
	}
	if len(producer.events) != 1 || producer.events[0].EventType != events.EventTypeMessageEdited {
		t.Fatalf("unexpected edit event: %+v", producer.events)
	}

	messages, err := svc.GetMessages(ctx, chatID, 10, nil)
	if err != nil {
		t.Fatalf("GetMessages returned error: %v", err)
	}
	if len(messages) != 1 || messages[0].ID != message.ID {
		t.Fatalf("unexpected messages: %+v", messages)
	}

	if err := svc.MarkAsRead(ctx, chatID, senderID); err != nil {
		t.Fatalf("MarkAsRead returned error: %v", err)
	}

	if err := svc.DeleteMessage(ctx, message.ID); err != nil {
		t.Fatalf("DeleteMessage returned error: %v", err)
	}
	if len(producer.events) != 2 || producer.events[1].EventType != events.EventTypeMessageDeleted {
		t.Fatalf("unexpected delete event: %+v", producer.events)
	}

	messages, err = svc.GetMessages(ctx, chatID, 10, nil)
	if err != nil {
		t.Fatalf("GetMessages returned error: %v", err)
	}
	if len(messages) != 0 {
		t.Fatalf("deleted message should be hidden, got %d messages", len(messages))
	}

	_, err = svc.EditMessage(ctx, uuid.New(), "x")
	if err == nil {
		t.Fatal("expected error for missing message")
	}
	if err := svc.DeleteMessage(ctx, uuid.New()); err == nil {
		t.Fatal("expected error for missing message delete")
	}
	if err := svc.MarkAsRead(ctx, uuid.New(), senderID); err == nil {
		t.Fatal("expected error for missing chat")
	}
	if _, err := svc.GetMessages(ctx, uuid.New(), 10, nil); err == nil {
		t.Fatal("expected error for missing chat messages")
	}
}

func TestForwardPinUnreadTypingAndReactions(t *testing.T) {
	svc, producer := newTestService()
	ctx := context.Background()
	aliceID := uuid.New()
	bobID := uuid.New()
	carolID := uuid.New()
	chatID := createChatWithMembers(t, svc, aliceID, bobID)
	targetChatID := createChatWithMembers(t, svc, aliceID, carolID)

	text := "source"
	source := createMessage(t, svc, chatID, aliceID, text)
	source.CreatedAt = time.Now().Add(-4 * time.Hour)

	if err := svc.attachmentRepo.Create(ctx, &model.Attachment{MessageID: source.ID, URL: "file", Type: model.AttachmentFile, FileName: "x", FileSize: 7}); err != nil {
		t.Fatalf("Create attachment returned error: %v", err)
	}

	forwarded, err := svc.ForwardMessage(ctx, source.ID, targetChatID, aliceID)
	if err != nil {
		t.Fatalf("ForwardMessage returned error: %v", err)
	}
	if forwarded == nil || forwarded.ChatID != targetChatID {
		t.Fatalf("unexpected forwarded message: %+v", forwarded)
	}
	forwardedAttachments, err := svc.attachmentRepo.GetByMessage(ctx, forwarded.ID)
	if err != nil {
		t.Fatalf("GetByMessage returned error: %v", err)
	}
	if len(forwardedAttachments) != 1 {
		t.Fatalf("unexpected forwarded attachments count: %d", len(forwardedAttachments))
	}
	if len(producer.events) != 1 || producer.events[0].EventType != events.EventTypeMessageSent {
		t.Fatalf("unexpected forward event: %+v", producer.events)
	}

	_, err = svc.ForwardMessage(ctx, source.ID, targetChatID, bobID)
	if !errors.Is(err, appErrors.ErrIsNotMember) {
		t.Fatalf("expected ErrIsNotMember, got %v", err)
	}

	if err := svc.PinMessage(ctx, source.ID); err != nil {
		t.Fatalf("PinMessage returned error: %v", err)
	}
	pinned, err := svc.GetPinnedMessages(ctx, chatID)
	if err != nil {
		t.Fatalf("GetPinnedMessages returned error: %v", err)
	}
	if len(pinned) != 1 || pinned[0].ID != source.ID {
		t.Fatalf("unexpected pinned messages: %+v", pinned)
	}
	if err := svc.UnpinMessage(ctx, source.ID); err != nil {
		t.Fatalf("UnpinMessage returned error: %v", err)
	}
	pinned, err = svc.GetPinnedMessages(ctx, chatID)
	if err != nil {
		t.Fatalf("GetPinnedMessages returned error: %v", err)
	}
	if len(pinned) != 0 {
		t.Fatalf("unexpected pinned messages after unpin: %+v", pinned)
	}
	if err := svc.PinMessage(ctx, uuid.New()); err == nil {
		t.Fatal("expected error for pinning missing message")
	}

	memberRepo := svc.membersRepo.(*memory.MembersRepository)
	reader := &model.ChatMember{ChatID: chatID, UserID: carolID, Role: model.Member, JoinedAt: time.Now()}
	if err := memberRepo.AddMember(ctx, reader); err != nil {
		t.Fatalf("AddMember returned error: %v", err)
	}
	first := createMessage(t, svc, chatID, aliceID, "one")
	second := createMessage(t, svc, chatID, aliceID, "two")
	third := createMessage(t, svc, chatID, aliceID, "three")
	first.CreatedAt = time.Now().Add(-3 * time.Hour)
	second.CreatedAt = time.Now().Add(-2 * time.Hour)
	third.CreatedAt = time.Now().Add(-1 * time.Hour)
	reader.LastReadMessageID = second.ID

	unreadInChat, err := svc.GetChatUnreadCount(ctx, chatID, carolID)
	if err != nil {
		t.Fatalf("GetChatUnreadCount returned error: %v", err)
	}
	if unreadInChat != 1 {
		t.Fatalf("unexpected unreadInChat: %d", unreadInChat)
	}

	otherChatID := createChatWithMembers(t, svc, carolID)
	otherReader := &model.ChatMember{ChatID: otherChatID, UserID: carolID, Role: model.Member, JoinedAt: time.Now()}
	if err := memberRepo.AddMember(ctx, otherReader); err != nil && !errors.Is(err, repository.ErrMemberExists) {
		t.Fatalf("AddMember returned error: %v", err)
	}
	otherMsg := createMessage(t, svc, otherChatID, carolID, "extra")
	otherMsg.CreatedAt = time.Now().Add(-30 * time.Minute)
	otherReader.LastReadMessageID = uuid.Nil

	totalUnread, err := svc.GetUnreadCount(ctx, carolID)
	if err != nil {
		t.Fatalf("GetUnreadCount returned error: %v", err)
	}
	if totalUnread < 2 {
		t.Fatalf("unexpected totalUnread: %d", totalUnread)
	}

	if err := svc.SendTyping(ctx, chatID, carolID); err != nil {
		t.Fatalf("SendTyping returned error: %v", err)
	}
	if err := svc.SendTyping(ctx, chatID, uuid.New()); !errors.Is(err, appErrors.ErrMemberNotFound) {
		t.Fatalf("expected ErrMemberNotFound, got %v", err)
	}

	if err := svc.AddReaction(ctx, source.ID, bobID, "🔥"); err != nil {
		t.Fatalf("AddReaction returned error: %v", err)
	}
	if len(producer.events) < 2 || producer.events[len(producer.events)-1].EventType != events.EventTypeReactionAdded {
		t.Fatalf("unexpected reaction add event: %+v", producer.events)
	}
	if err := svc.RemoveReaction(ctx, source.ID, bobID); err != nil {
		t.Fatalf("RemoveReaction returned error: %v", err)
	}
	if producer.events[len(producer.events)-1].EventType != events.EventTypeReactionRemoved {
		t.Fatalf("unexpected reaction remove event: %+v", producer.events)
	}
	if err := svc.AddReaction(ctx, uuid.New(), bobID, "🔥"); err == nil {
		t.Fatal("expected error for missing message reaction")
	}
	if err := svc.RemoveReaction(ctx, uuid.New(), bobID); err == nil {
		t.Fatal("expected error for missing message reaction removal")
	}
}

func TestHelpersAndConstructor(t *testing.T) {
	svc, _ := newTestService()
	constructed := NewService(
		svc.chatRepo,
		svc.membersRepo,
		svc.messageRepo,
		svc.attachmentRepo,
		svc.reactionRepo,
		svc.pinnedRepo,
		svc.typingRepo,
		nil,
		nil,
		nil,
	)
	if constructed == nil {
		t.Fatal("NewService returned nil")
	}

	if messageText(nil) != "" {
		t.Fatal("messageText should return empty string for nil input")
	}
	short := "short"
	if messageText(&short) != short {
		t.Fatal("messageText should keep short strings")
	}
	long := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz0123456789"
	result := messageText(&long)
	if len(result) != 103 || result[len(result)-3:] != "..." {
		t.Fatalf("unexpected truncated text: %q", result)
	}

	users := []uuid.UUID{uuid.New(), uuid.New()}
	ids := getMemberIDs([]*model.ChatMember{{UserID: users[0]}, {UserID: users[1]}})
	if len(ids) != 2 || ids[0] != users[0] || ids[1] != users[1] {
		t.Fatalf("unexpected ids: %+v", ids)
	}

	chatID := createChatWithMembers(t, svc, users[0], users[1])
	receiverID, err := svc.findChatReceiver(context.Background(), chatID, users[0])
	if err != nil {
		t.Fatalf("findChatReceiver returned error: %v", err)
	}
	if receiverID != users[1] {
		t.Fatalf("unexpected receiver: %s", receiverID)
	}

	soloService, _ := newTestService()
	soloChatID := createChatWithMembers(t, soloService, users[0])
	_, err = soloService.findChatReceiver(context.Background(), soloChatID, users[0])
	if !errors.Is(err, appErrors.ErrNoMembers) {
		t.Fatalf("expected ErrNoMembers, got %v", err)
	}

	soloService.producer = nil
	if err := soloService.publishEvent(events.Event{}); err != nil {
		t.Fatalf("publishEvent with nil producer should not fail: %v", err)
	}

	wrappedService, _ := newTestService()
	wrappedChatID := createChatWithMembers(t, wrappedService, uuid.New())
	wrappedService.messageRepo = nilMessageRepo{wrappedService.messageRepo}
	_, err = wrappedService.GetMessages(context.Background(), wrappedChatID, 10, nil)
	if !errors.Is(err, appErrors.ErrMessageNotFound) {
		t.Fatalf("expected ErrMessageNotFound, got %v", err)
	}
}
