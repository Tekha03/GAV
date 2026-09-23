package grpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	authv1 "api/gen/auth/v1"
	notificationv1 "api/gen/notification/v1"
	socialv1 "api/gen/social/v1"
	"social_network/internal/app"
	"social_network/internal/auth"
	"social_network/internal/config"
	"social_network/internal/profile"
	"social_network/internal/testdb"
	grpctransport "social_network/transport/grpc"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestMessengerSocialNetworkContract(t *testing.T) {
	t.Setenv("ENV", "production")
	const secret = "contract-test-secret"
	cfg := &config.Config{
		HTTP:    config.HTTPConfig{Port: "8080"},
		GRPC:    config.GRPCConfig{Addr: ":9000"},
		DB:      config.DBConfig{PostgresDSN: testdb.DSN(t)},
		JWT:     config.JWTConfig{Secret: secret, TTL: time.Hour},
		Storage: config.StorageConfig{Type: "local", LocalPath: t.TempDir()},
	}

	application, err := app.NewApp(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = application.Shutdown(context.Background()) })

	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	grpctransport.Register(server, application.Services.User, application.Services.Profile, application.Services.Auth, application.Services.Notification, cfg.JWT.TTL)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	authClient := authv1.NewAuthServiceClient(conn)
	registered, err := authClient.Register(context.Background(), &authv1.RegisterRequest{
		Email: "contract@example.com", Password: "password-123", Username: "contract-user",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	claims, err := auth.ParseToken(registered.GetTokens().GetAccessToken(), auth.JWTConfig{Secret: []byte(secret), TTL: time.Hour})
	if err != nil {
		t.Fatalf("parse registered JWT: %v", err)
	}

	_, err = application.Services.Profile.Create(context.Background(), claims.UserID, profile.CreateProfileInput{Username: "contract-user", Name: "Contract"})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}

	loggedIn, err := authClient.Login(context.Background(), &authv1.LoginRequest{Email: "contract@example.com", Password: "password-123"})
	if err != nil || loggedIn.GetTokens().GetAccessToken() == "" {
		t.Fatalf("login: tokens=%v err=%v", loggedIn.GetTokens(), err)
	}

	profileResponse, err := socialv1.NewSocialServiceClient(conn).GetProfile(context.Background(), &socialv1.GetProfileRequest{UserId: claims.UserID.String()})
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if profileResponse.GetUsername() != "contract-user" {
		t.Fatalf("username = %q", profileResponse.GetUsername())
	}

	receiverID := uuid.New()
	chatID := uuid.New()
	response, err := notificationv1.NewNotificationServiceClient(conn).SendNotification(context.Background(), &notificationv1.NotificationRequest{
		UserId: receiverID[:], EventType: "new_message",
		Data: map[string]string{"sender_id": claims.UserID.String(), "chat_id": chatID.String()},
	})
	if err != nil || !response.GetSuccess() {
		t.Fatalf("send notification: response=%v err=%v", response, err)
	}
	notifications, err := application.Services.Notification.GetInAppNotifications(context.Background(), receiverID)
	if err != nil || len(notifications) != 1 {
		t.Fatalf("stored notifications=%d err=%v", len(notifications), err)
	}
	if notifications[0].UserID != receiverID || notifications[0].FromUserID != claims.UserID || notifications[0].EntityID != chatID {
		t.Fatalf("notification routed incorrectly: %+v", notifications[0])
	}
}
