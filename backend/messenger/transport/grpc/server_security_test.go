package grpc

import (
	"context"
	"testing"

	pb "api/gen/chat/v1"
	"messenger/internal/model"
	"messenger/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type securityCaptureService struct {
	service.Service
	privateUsers      [2]uuid.UUID
	groupCreator      uuid.UUID
	addedChatID       uuid.UUID
	addedUserID       uuid.UUID
	addRequesterID    uuid.UUID
	messageRequester  uuid.UUID
	messageInput      model.SendMessageInput
	reactionUser      uuid.UUID
	reactionRequester uuid.UUID
}

func (s *securityCaptureService) CreatePrivateChat(_ context.Context, userID1, userID2 uuid.UUID) (*model.Chat, error) {
	s.privateUsers = [2]uuid.UUID{userID1, userID2}
	return &model.Chat{ID: uuid.New()}, nil
}

func (s *securityCaptureService) CreateGroupChat(_ context.Context, _ string, creatorID uuid.UUID, _ []uuid.UUID) (*model.Chat, error) {
	s.groupCreator = creatorID
	return &model.Chat{ID: uuid.New()}, nil
}

func (s *securityCaptureService) AddMember(_ context.Context, chatID, userID, requesterID uuid.UUID) error {
	s.addedChatID = chatID
	s.addedUserID = userID
	s.addRequesterID = requesterID
	return nil
}

func (s *securityCaptureService) SendMessage(_ context.Context, requesterID uuid.UUID, input model.SendMessageInput) (*model.Message, error) {
	s.messageRequester = requesterID
	s.messageInput = input
	return &model.Message{ID: uuid.New()}, nil
}

func (s *securityCaptureService) AddReaction(_ context.Context, _ uuid.UUID, userID, requesterID uuid.UUID, _ string) error {
	s.reactionUser = userID
	s.reactionRequester = requesterID
	return nil
}

func authenticatedContext(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), userIDContextKey, userID)
}

func TestCreatePrivateChatRejectsUnrelatedRequester(t *testing.T) {
	capture := &securityCaptureService{}
	server := NewServer(capture)

	_, err := server.CreatePrivateChat(authenticatedContext(uuid.New()), &pb.CreatePrivateChatRequest{
		UserId_1: uuid.New().String(),
		UserId_2: uuid.New().String(),
	})

	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("code = %s, want %s", status.Code(err), codes.PermissionDenied)
	}
	if capture.privateUsers != [2]uuid.UUID{} {
		t.Fatal("service must not be called for an unrelated requester")
	}
}

func TestSendMessageUsesAuthenticatedUserAsSender(t *testing.T) {
	requesterID := uuid.New()
	spoofedSenderID := uuid.New()
	capture := &securityCaptureService{}
	server := NewServer(capture)

	_, err := server.SendMessage(authenticatedContext(requesterID), &pb.SendMessageRequest{
		ChatId:   uuid.New().String(),
		SenderId: spoofedSenderID.String(),
		Text:     "hello",
	})
	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	if capture.messageRequester != requesterID || capture.messageInput.SenderID != requesterID {
		t.Fatalf("requester = %s, sender = %s, want %s", capture.messageRequester, capture.messageInput.SenderID, requesterID)
	}
}

func TestActorIDsComeFromAuthenticatedContext(t *testing.T) {
	requesterID := uuid.New()
	targetUserID := uuid.New()
	chatID := uuid.New()
	messageID := uuid.New()
	capture := &securityCaptureService{}
	server := NewServer(capture)
	ctx := authenticatedContext(requesterID)

	if _, err := server.CreateGroupChat(ctx, &pb.CreateGroupChatRequest{
		CreatorId: uuid.New().String(),
	}); err != nil {
		t.Fatalf("CreateGroupChat() error = %v", err)
	}
	if capture.groupCreator != requesterID {
		t.Fatalf("creator = %s, want %s", capture.groupCreator, requesterID)
	}

	if _, err := server.AddMember(ctx, &pb.AddMemberRequest{
		ChatId: chatID.String(),
		UserId: targetUserID.String(),
	}); err != nil {
		t.Fatalf("AddMember() error = %v", err)
	}
	if capture.addedChatID != chatID || capture.addedUserID != targetUserID || capture.addRequesterID != requesterID {
		t.Fatalf("AddMember IDs = (%s, %s, %s)", capture.addedChatID, capture.addedUserID, capture.addRequesterID)
	}

	if _, err := server.AddReaction(ctx, &pb.AddReactionRequest{
		MessageId: messageID.String(),
		UserId:    uuid.New().String(),
		Emoji:     "👍",
	}); err != nil {
		t.Fatalf("AddReaction() error = %v", err)
	}
	if capture.reactionUser != requesterID || capture.reactionRequester != requesterID {
		t.Fatalf("reaction IDs = (%s, %s), want %s", capture.reactionUser, capture.reactionRequester, requesterID)
	}
}
