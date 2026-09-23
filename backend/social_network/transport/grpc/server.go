package grpc

import (
	"context"
	"time"

	authv1 "api/gen/auth/v1"
	notificationv1 "api/gen/notification/v1"
	socialv1 "api/gen/social/v1"
	"social_network/internal/auth"
	"social_network/internal/notification"
	"social_network/internal/profile"
	"social_network/internal/user"

	"github.com/google/uuid"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SocialServer struct {
	socialv1.UnimplementedSocialServiceServer
	users    user.UserService
	profiles profile.ProfileService
}

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	auth auth.AuthService
	ttl  time.Duration
}

type NotificationServer struct {
	notificationv1.UnimplementedNotificationServiceServer
	notifications notification.NotificationService
}

func Register(server *grpcpkg.Server, users user.UserService, profiles profile.ProfileService, authService auth.AuthService, notifications notification.NotificationService, ttl time.Duration) {
	socialv1.RegisterSocialServiceServer(server, &SocialServer{users: users, profiles: profiles})
	authv1.RegisterAuthServiceServer(server, &AuthServer{auth: authService, ttl: ttl})
	notificationv1.RegisterNotificationServiceServer(server, &NotificationServer{notifications: notifications})
}

func (s *SocialServer) GetUser(ctx context.Context, req *socialv1.GetUserRequest) (*socialv1.User, error) {
	userID, err := parseUUID(req.GetUserId(), "user_id")
	if err != nil {
		return nil, err
	}
	account, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	result := &socialv1.User{
		Id: account.ID.String(), Email: account.Email, Role: account.Role,
		CreatedAt: timestamppb.New(account.CreatedAt), UpdatedAt: timestamppb.New(account.UpdatedAt),
	}
	if p, profileErr := s.profiles.GetByUserID(ctx, userID); profileErr == nil {
		result.Profile = profileToProto(p)
	}
	return result, nil
}

func (s *SocialServer) GetProfile(ctx context.Context, req *socialv1.GetProfileRequest) (*socialv1.UserProfile, error) {
	userID, err := parseUUID(req.GetUserId(), "user_id")
	if err != nil {
		return nil, err
	}
	p, err := s.profiles.GetByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}
	return profileToProto(p), nil
}

func (s *AuthServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	tokens, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &authv1.RegisterResponse{Tokens: tokensToProto(tokens, s.ttl)}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	tokens, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &authv1.LoginResponse{Tokens: tokensToProto(tokens, s.ttl)}, nil
}

func (s *AuthServer) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	tokens, err := s.auth.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}
	return &authv1.RefreshResponse{Tokens: tokensToProto(tokens, s.ttl)}, nil
}

func (s *AuthServer) Me(context.Context, *emptypb.Empty) (*authv1.MeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "use the public HTTP auth endpoint")
}

func (s *AuthServer) Logout(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "use the public HTTP auth endpoint")
}

func (s *NotificationServer) SendNotification(ctx context.Context, req *notificationv1.NotificationRequest) (*notificationv1.NotificationResponse, error) {
	receiverID, err := uuid.FromBytes(req.GetUserId())
	if err != nil || receiverID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	if req.GetEventType() != "new_message" {
		return nil, status.Error(codes.InvalidArgument, "unsupported notification event_type")
	}
	senderID, err := parseUUID(req.GetData()["sender_id"], "sender_id")
	if err != nil {
		return nil, err
	}
	chatID, err := parseUUID(req.GetData()["chat_id"], "chat_id")
	if err != nil {
		return nil, err
	}
	if err := s.notifications.NotifyNewMessage(ctx, receiverID, senderID, chatID); err != nil {
		return nil, status.Error(codes.Internal, "failed to persist notification")
	}
	return &notificationv1.NotificationResponse{Success: true}, nil
}

func parseUUID(raw, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "invalid %s", field)
	}
	return id, nil
}

func profileToProto(p *profile.UserProfile) *socialv1.UserProfile {
	return &socialv1.UserProfile{
		UserId: p.UserID.String(), Name: p.Name, Surname: p.Surname, Username: p.Username,
		ProfilePhotoUrl: p.ProfilePhotoUrl, Bio: p.Bio, Address: p.Address, BirthDate: p.BirthDate,
	}
}

func tokensToProto(tokens *auth.AuthTokens, ttl time.Duration) *authv1.AuthTokens {
	return &authv1.AuthTokens{
		AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken,
		AccessExpiresIn: int32(ttl.Seconds()),
	}
}
