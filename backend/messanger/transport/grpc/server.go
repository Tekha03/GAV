package grpc

import (
	"context"
	"time"

	pb "api/chat_gen/chat"
	"messanger/internal/model"
	"messanger/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
    chat, err := s.service.CreatePrivateChat(ctx, 
        uuid.MustParse(req.UserId_1), 
        uuid.MustParse(req.UserId_2),
    )
    if err != nil {
        return nil, status.Errorf(codes.Internal, "create private chat: %v", err)
    }
    return &pb.CreatePrivateChatResponse{ChatId: chat.ID.String()}, nil
}

func (s *Server) CreateGroupChat(ctx context.Context, req *pb.CreateGroupChatRequest) (*pb.CreateGroupChatResponse, error) {
    memberIDs := make([]uuid.UUID, len(req.MemberIds))
    for i, id := range req.MemberIds {
        memberIDs[i] = uuid.MustParse(id)
    }
    
    chat, err := s.service.CreateGroupChat(ctx, req.Title, uuid.MustParse(req.CreatorId), memberIDs)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "create group chat: %v", err)
    }
    return &pb.CreateGroupChatResponse{ChatId: chat.ID.String()}, nil
}

func (s *Server) GetChat(ctx context.Context, req *pb.GetChatRequest) (*pb.GetChatResponse, error) {
    chat, err := s.service.GetChatByID(ctx, uuid.MustParse(req.ChatId))
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "chat not found: %v", err)
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
    err := s.service.AddMember(ctx, uuid.MustParse(req.UserId), uuid.MustParse(req.ChatId))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "add member: %v", err)
    }
    return &pb.AddMemberResponse{Success: true}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
    var replyTo *uuid.UUID
    if req.ReplyToId != "" {
        id := uuid.MustParse(req.ReplyToId)
        replyTo = &id
    }
    
    input := model.SendMessageInput{
        ChatID:    uuid.MustParse(req.ChatId),
        SenderID:  uuid.MustParse(req.SenderId),
        Text:      &req.Text,
        ReplyToID: replyTo,
    }
    
    msg, err := s.service.SendMessage(ctx, input)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "send message: %v", err)
    }
    
    return &pb.SendMessageResponse{
        MessageId: msg.ID.String(),
        Success:   true,
    }, nil
}

func (s *Server) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
    var cursor *uuid.UUID
    if req.CursorId != "" {
        id := uuid.MustParse(req.CursorId)
        cursor = &id
    }
    
    messages, err := s.service.GetMessages(ctx, uuid.MustParse(req.ChatId), int(req.Limit), cursor)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "get messages: %v", err)
    }
    
    pbMessages := make([]*pb.Message, len(messages))
    for i, msg := range messages {
        pbMsg := &pb.Message{
            Id:          msg.ID.String(),
            ChatId:      msg.ChatID.String(),
            SenderId:    msg.SenderID.String(),
            Text:        *msg.Text,
            ReplyToId:   uuidOrEmpty(msg.ReplyToID),
            CreatedAt:   timestamppb.New(msg.CreatedAt),
            EditedAt:    timestamppbOrNil(msg.EditedAt),
            DeletedAt:   timestamppbOrNil(msg.DeletedAt),
            ReadAt:      timestamppbOrNil(msg.ReadAt),
        }
        pbMessages[i] = pbMsg
    }
    
    return &pb.GetMessagesResponse{Messages: pbMessages}, nil
}

func (s *Server) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
    _, err := s.service.EditMessage(ctx, uuid.MustParse(req.MessageId), req.NewText)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "edit message: %v", err)
    }
    return &pb.EditMessageResponse{Success: true}, nil
}

func (s *Server) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
    err := s.service.DeleteMessage(ctx, uuid.MustParse(req.MessageId))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "delete message: %v", err)
    }
    return &pb.DeleteMessageResponse{Success: true}, nil
}

func (s *Server) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
    err := s.service.MarkAsRead(ctx, uuid.MustParse(req.ChatId), uuid.MustParse(req.UserId))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "mark as read: %v", err)
    }
    return &pb.MarkAsReadResponse{Success: true}, nil
}

func (s *Server) AddReaction(ctx context.Context, req *pb.AddReactionRequest) (*pb.AddReactionResponse, error) {
    err := s.service.AddReaction(ctx, uuid.MustParse(req.MessageId), uuid.MustParse(req.UserId), req.Emoji)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "add reaction: %v", err)
    }
    return &pb.AddReactionResponse{Success: true}, nil
}

func (s *Server) RemoveReaction(ctx context.Context, req *pb.RemoveReactionRequest) (*pb.RemoveReactionResponse, error) {
    err := s.service.RemoveReaction(ctx, uuid.MustParse(req.MessageId), uuid.MustParse(req.UserId))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "remove reaction: %v", err)
    }
    return &pb.RemoveReactionResponse{Success: true}, nil
}

func (s *Server) PinMessage(ctx context.Context, req *pb.PinMessageRequest) (*pb.PinMessageResponse, error) {
    err := s.service.PinMessage(ctx, uuid.MustParse(req.MessageId))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "pin message: %v", err)
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