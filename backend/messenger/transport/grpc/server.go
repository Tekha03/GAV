package grpc

import (
	"context"
	"time"

	pb "api/gen/chat/v1"
	"messenger/internal/model"
	"messenger/internal/service"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedChatServiceServer
	service service.Service
}

func NewServer(service service.Service) *Server {
	return &Server{service: service}
}

func (s *Server) CreatePrivateChat(ctx context.Context, req *pb.CreatePrivateChatRequest) (*pb.CreatePrivateChatResponse, error) {
	userID1, err := parseUUID(req.UserId_1, "user_id_1")
	if err != nil {
		return nil, toGRPCError(err)
	}
	userID2, err := parseUUID(req.UserId_2, "user_id_2")
	if err != nil {
		return nil, toGRPCError(err)
	}
	chat, err := s.service.CreatePrivateChat(ctx, userID1, userID2)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.CreatePrivateChatResponse{ChatId: chat.ID.String()}, nil
}

func (s *Server) CreateGroupChat(ctx context.Context, req *pb.CreateGroupChatRequest) (*pb.CreateGroupChatResponse, error) {
	memberIDs := make([]uuid.UUID, len(req.MemberIds))
	for i, id := range req.MemberIds {
		memberID, err := parseUUID(id, "member_ids")
		if err != nil {
			return nil, toGRPCError(err)
		}
		memberIDs[i] = memberID
	}
	creatorID, err := parseUUID(req.CreatorId, "creator_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	chat, err := s.service.CreateGroupChat(ctx, req.Title, creatorID, memberIDs)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.CreateGroupChatResponse{ChatId: chat.ID.String()}, nil
}

func (s *Server) GetChat(ctx context.Context, req *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	requesterID, err := requireCurrentUserID(CurrentUserID(ctx))
	if err != nil {
		return nil, toGRPCError(err)
	}
	chatID, err := parseUUID(req.ChatId, "chat_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	chat, err := s.service.GetChatByID(ctx, chatID, requesterID)
	if err != nil {
		return nil, toGRPCError(err)
	}

	pbChat := &pb.Chat{
		Id:        chat.ID.String(),
		IsGroup:   chat.IsGroup,
		Title:     chat.Title,
		PhotoUrl:  chat.PhotoURL,
		CreatedAt: timestamppb.New(chat.CreatedAt),
	}
	return &pb.GetChatResponse{Chat: pbChat}, nil
}

func (s *Server) AddMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.AddMemberResponse, error) {
	chatID, err := parseUUID(req.ChatId, "chat_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	err = s.service.AddMember(ctx, chatID, userID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.AddMemberResponse{Success: true}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	var replyTo *uuid.UUID
	if req.ReplyToId != "" {
		id, err := parseUUID(req.ReplyToId, "reply_to_id")
		if err != nil {
			return nil, toGRPCError(err)
		}
		replyTo = &id
	}
	chatID, err := parseUUID(req.ChatId, "chat_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	senderID, err := parseUUID(req.SenderId, "sender_id")
	if err != nil {
		return nil, toGRPCError(err)
	}

	input := model.SendMessageInput{
		ChatID:    chatID,
		SenderID:  senderID,
		Text:      &req.Text,
		ReplyToID: replyTo,
	}

	msg, err := s.service.SendMessage(ctx, input.SenderID, input)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &pb.SendMessageResponse{
		MessageId: msg.ID.String(),
		Success:   true,
	}, nil
}

func (s *Server) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	var cursor *uuid.UUID
	if req.CursorId != "" {
		id, err := parseUUID(req.CursorId, "cursor_id")
		if err != nil {
			return nil, toGRPCError(err)
		}
		cursor = &id
	}
	requesterID, err := requireCurrentUserID(CurrentUserID(ctx))
	if err != nil {
		return nil, toGRPCError(err)
	}
	chatID, err := parseUUID(req.ChatId, "chat_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	messages, err := s.service.GetMessages(ctx, chatID, requesterID, int(req.Limit), cursor)
	if err != nil {
		return nil, toGRPCError(err)
	}

	pbMessages := make([]*pb.Message, len(messages))
	for i, msg := range messages {
		text := ""
		if msg.Text != nil {
			text = *msg.Text
		}
		pbMsg := &pb.Message{
			Id:        msg.ID.String(),
			ChatId:    msg.ChatID.String(),
			SenderId:  msg.SenderID.String(),
			Text:      text,
			ReplyToId: uuidOrEmpty(msg.ReplyToID),
			CreatedAt: timestamppb.New(msg.CreatedAt),
			EditedAt:  timestamppbOrNil(msg.EditedAt),
			DeletedAt: timestamppbOrNil(msg.DeletedAt),
			ReadAt:    timestamppbOrNil(msg.ReadAt),
		}
		for _, attachment := range msg.Attachments {
			pbMsg.Attachments = append(pbMsg.Attachments, &pb.Attachment{
				Id:        attachment.ID.String(),
				MessageId: attachment.MessageID.String(),
				Url:       attachment.URL,
				Type:      toPBAttachmentType(attachment.Type),
				FileName:  attachment.FileName,
				FileSize:  attachment.FileSize,
			})
		}
		pbMessages[i] = pbMsg
	}

	return &pb.GetMessagesResponse{Messages: pbMessages}, nil
}

func (s *Server) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
	requesterID, err := requireCurrentUserID(CurrentUserID(ctx))
	if err != nil {
		return nil, toGRPCError(err)
	}
	messageID, err := parseUUID(req.MessageId, "message_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	_, err = s.service.EditMessage(ctx, requesterID, messageID, req.NewText)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.EditMessageResponse{Success: true}, nil
}

func (s *Server) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	requesterID, err := requireCurrentUserID(CurrentUserID(ctx))
	if err != nil {
		return nil, toGRPCError(err)
	}
	messageID, err := parseUUID(req.MessageId, "message_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	err = s.service.DeleteMessage(ctx, requesterID, messageID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.DeleteMessageResponse{Success: true}, nil
}

func (s *Server) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	requesterID, err := requireCurrentUserID(CurrentUserID(ctx))
	if err != nil {
		return nil, toGRPCError(err)
	}
	chatID, err := parseUUID(req.ChatId, "chat_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	err = s.service.MarkAsRead(ctx, chatID, requesterID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.MarkAsReadResponse{Success: true}, nil
}

func (s *Server) AddReaction(ctx context.Context, req *pb.AddReactionRequest) (*pb.AddReactionResponse, error) {
	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	messageID, err := parseUUID(req.MessageId, "message_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	err = s.service.AddReaction(ctx, messageID, userID, userID, req.Emoji)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.AddReactionResponse{Success: true}, nil
}

func (s *Server) RemoveReaction(ctx context.Context, req *pb.RemoveReactionRequest) (*pb.RemoveReactionResponse, error) {
	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	messageID, err := parseUUID(req.MessageId, "message_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	err = s.service.RemoveReaction(ctx, messageID, userID, userID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.RemoveReactionResponse{Success: true}, nil
}

func (s *Server) PinMessage(ctx context.Context, req *pb.PinMessageRequest) (*pb.PinMessageResponse, error) {
	requesterID, err := requireCurrentUserID(CurrentUserID(ctx))
	if err != nil {
		return nil, toGRPCError(err)
	}
	messageID, err := parseUUID(req.MessageId, "message_id")
	if err != nil {
		return nil, toGRPCError(err)
	}
	err = s.service.PinMessage(ctx, requesterID, messageID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.PinMessageResponse{Success: true}, nil
}

func uuidOrEmpty(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

func timestamppbOrNil(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func toPBAttachmentType(attachmentType model.AttachmentType) pb.AttachmentType {
	switch attachmentType {
	case model.AttachmentImage:
		return pb.AttachmentType_ATTACHMENT_IMAGE
	case model.AttachmentVideo:
		return pb.AttachmentType_ATTACHMENT_VIDEO
	case model.AttachmentVoice:
		return pb.AttachmentType_ATTACHMENT_VOICE
	default:
		return pb.AttachmentType_ATTACHMENT_FILE
	}
}
