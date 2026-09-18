package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"messenger/internal/model"
	"messenger/internal/service"
	storagegorm "messenger/storage/gorm"
	"messenger/transport/websocket"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func e2ePostgres(t *testing.T) *gorm.DB {
	t.Helper()
	base := os.Getenv("GAV_TEST_POSTGRES_DSN")
	if base == "" {
		t.Skip("GAV_TEST_POSTGRES_DSN is not set")
	}
	admin, err := gorm.Open(postgres.Open(base), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	schema := "messenger_e2e_" + strings.ReplaceAll(uuid.NewString()[:8], "-", "")
	if err := admin.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schema)).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec(fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, schema)).Error; err != nil {
			t.Error(err)
		}
		connection, _ := admin.DB()
		_ = connection.Close()
	})
	u, err := url.Parse(base)
	if err != nil {
		t.Fatal(err)
	}
	query := u.Query()
	query.Set("search_path", schema)
	u.RawQuery = query.Encode()
	db, err := gorm.Open(postgres.Open(u.String()), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		connection, _ := db.DB()
		_ = connection.Close()
	})
	if err := db.AutoMigrate(&model.Chat{}, &model.ChatMember{}, &model.Message{}, &model.Attachment{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestMessengerPrivateDeliveryAndReconnectPostgres(t *testing.T) {
	db := e2ePostgres(t)
	repo := storagegorm.NewRepository(db)
	hub := websocket.NewHub()
	go hub.Run()
	chatService := service.NewService(
		repo, storagegorm.NewChatRepository(repo), storagegorm.NewChatMemberRepository(repo),
		storagegorm.NewMessageRepository(repo), storagegorm.NewAttachmentRepository(repo), nil,
		storagegorm.NewReactionRepository(repo), nil, nil, nil, nil, nil, hub,
	)
	const secret = "messenger-e2e-secret"
	server := httptest.NewServer(NewHTTPServer("", chatService, secret, hub).Handler)
	defer server.Close()

	sender, recipient, outsider := uuid.New(), uuid.New(), uuid.New()
	token := func(id uuid.UUID) string {
		t.Helper()
		value, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": id.String(), "role": "user", "exp": time.Now().Add(time.Hour).Unix(),
		}).SignedString([]byte(secret))
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	request := func(method, path string, id uuid.UUID, body string) *http.Response {
		t.Helper()
		req, err := http.NewRequest(method, server.URL+path, bytes.NewBufferString(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+token(id))
		req.Header.Set("Content-Type", "application/json")
		response, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = response.Body.Close() })
		return response
	}
	chatResponse := request(http.MethodPost, "/api/v1/chats/private", sender,
		fmt.Sprintf(`{"user_id_1":%q,"user_id_2":%q}`, sender, recipient))
	if chatResponse.StatusCode != http.StatusCreated {
		t.Fatalf("create chat status = %d", chatResponse.StatusCode)
	}
	var created struct {
		Chat struct {
			ID uuid.UUID `json:"id"`
		} `json:"chat"`
	}
	if err := json.NewDecoder(chatResponse.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	chatID := created.Chat.ID
	path := "/api/v1/chats/" + chatID.String()
	otherChatResponse := request(http.MethodPost, "/api/v1/chats/private", sender,
		fmt.Sprintf(`{"user_id_1":%q,"user_id_2":%q}`, sender, outsider))
	if otherChatResponse.StatusCode != http.StatusCreated {
		t.Fatalf("create unrelated chat status = %d", otherChatResponse.StatusCode)
	}
	var otherCreated struct {
		Chat struct {
			ID uuid.UUID `json:"id"`
		} `json:"chat"`
	}
	if err := json.NewDecoder(otherChatResponse.Body).Decode(&otherCreated); err != nil {
		t.Fatal(err)
	}
	if got := request(http.MethodGet, path+"/messages", outsider, "").StatusCode; got != http.StatusForbidden {
		t.Fatalf("outsider history status = %d, want 403", got)
	}
	if got := request(http.MethodPost, path+"/messages", outsider,
		fmt.Sprintf(`{"sender_id":%q,"text":"intrusion"}`, outsider)).StatusCode; got != http.StatusForbidden {
		t.Fatalf("outsider send status = %d, want 403", got)
	}
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/chats/" + chatID.String()
	dial := func(id uuid.UUID) (*ws.Conn, *http.Response, error) {
		header := http.Header{"Authorization": {"Bearer " + token(id)}}
		return ws.DefaultDialer.Dial(wsURL, header)
	}
	if conn, response, err := dial(outsider); err == nil {
		_ = conn.Close()
		t.Fatal("outsider WebSocket connected")
	} else if response == nil || response.StatusCode != http.StatusForbidden {
		t.Fatalf("outsider WebSocket status = %v, error = %v", response, err)
	}
	recipientConn, response, err := dial(recipient)
	if err != nil {
		t.Fatalf("recipient WebSocket: %v, response: %v", err, response)
	}
	defer recipientConn.Close()
	otherURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/chats/" + otherCreated.Chat.ID.String()
	otherHeader := http.Header{"Authorization": {"Bearer " + token(outsider)}}
	otherConn, response, err := ws.DefaultDialer.Dial(otherURL, otherHeader)
	if err != nil {
		t.Fatalf("unrelated chat WebSocket: %v, response: %v", err, response)
	}
	defer otherConn.Close()

	send := func(text string) uuid.UUID {
		t.Helper()
		response := request(http.MethodPost, path+"/messages", sender,
			fmt.Sprintf(`{"sender_id":%q,"text":%q}`, sender, text))
		if response.StatusCode != http.StatusCreated {
			t.Fatalf("send status = %d", response.StatusCode)
		}
		var payload struct {
			Message struct {
				ID uuid.UUID `json:"id"`
			} `json:"message"`
		}
		if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		return payload.Message.ID
	}
	firstID := send("first message")
	if err := recipientConn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var event struct {
		Type    string    `json:"type"`
		ChatID  uuid.UUID `json:"chat_id"`
		Message struct {
			ID       uuid.UUID `json:"id"`
			SenderID uuid.UUID `json:"sender_id"`
		} `json:"message"`
	}
	if err := recipientConn.ReadJSON(&event); err != nil {
		t.Fatal(err)
	}
	if event.Type != "message.sent" || event.ChatID != chatID || event.Message.ID != firstID || event.Message.SenderID != sender {
		t.Fatalf("wrong recipient event: %+v", event)
	}
	_ = otherConn.SetReadDeadline(time.Now().Add(150 * time.Millisecond))
	if _, _, err := otherConn.ReadMessage(); err == nil {
		t.Fatal("message from private chat leaked to unrelated chat")
	} else if timeout, ok := err.(net.Error); !ok || !timeout.Timeout() {
		t.Fatalf("unrelated chat WebSocket read error = %v, want timeout", err)
	}
	_ = recipientConn.Close()
	secondID := send("while offline")
	if secondID == firstID {
		t.Fatal("duplicate message ID")
	}
	reconnected, response, err := dial(recipient)
	if err != nil {
		t.Fatalf("reconnect WebSocket: %v, response: %v", err, response)
	}
	defer reconnected.Close()
	history := request(http.MethodGet, path+"/messages", recipient, "")
	if history.StatusCode != http.StatusOK {
		t.Fatalf("recipient history status = %d", history.StatusCode)
	}
	var messages struct {
		Messages []struct {
			ID uuid.UUID `json:"id"`
		} `json:"messages"`
	}
	if err := json.NewDecoder(history.Body).Decode(&messages); err != nil {
		t.Fatal(err)
	}
	if len(messages.Messages) != 2 || messages.Messages[0].ID != secondID || messages.Messages[1].ID != firstID {
		t.Fatalf("reconnected history = %+v", messages.Messages)
	}
	var count int64
	if err := db.Model(&model.Message{}).Where("chat_id = ?", chatID).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("stored message count = %d, want 2", count)
	}
	_ = reconnected.SetReadDeadline(time.Now().Add(150 * time.Millisecond))
	if _, _, err := reconnected.ReadMessage(); err == nil {
		t.Fatal("reconnect replayed an old WebSocket event")
	} else if timeout, ok := err.(net.Error); !ok || !timeout.Timeout() {
		t.Fatalf("reconnected WebSocket read error = %v, want timeout", err)
	}
}
