package handlers

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"social_network/internal/notification"
	"social_network/transport/http/middleware"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func wsURL(serverURL string) string {
	return "ws" + serverURL[len("http"):]
}

func requireLocalListener(t *testing.T) {
	t.Helper()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Skipf("local listener unavailable: %v", err)
	}
	listener.Close()
}

func TestNotificationHandler_ServeWS(t *testing.T) {
	userID := uuid.New()

	t.Run("unauthorized", func(t *testing.T) {
		hub := notification.NewHub()
		go hub.Run()

		h, _ := NewNotificationHandler(hub)

		req := httptest.NewRequest("GET", "/ws", nil)
		w := httptest.NewRecorder()

		h.ServeWS(w, req)

		assert.NotEqual(t, http.StatusSwitchingProtocols, w.Code)
	})

	t.Run("successful connection", func(t *testing.T) {
		requireLocalListener(t)

		hub := notification.NewHub()
		go hub.Run()

		h, _ := NewNotificationHandler(hub)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
			h.ServeWS(w, r.WithContext(ctx))
		}))
		defer server.Close()

		conn, _, err := websocket.DefaultDialer.Dial(wsURL(server.URL), nil)
		assert.NoError(t, err)
		defer conn.Close()
	})

	t.Run("broadcast message received", func(t *testing.T) {
		requireLocalListener(t)

		hub := notification.NewHub()
		go hub.Run()

		h, _ := NewNotificationHandler(hub)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
			h.ServeWS(w, r.WithContext(ctx))
		}))
		defer server.Close()

		conn, _, err := websocket.DefaultDialer.Dial(wsURL(server.URL), nil)
		assert.NoError(t, err)
		defer conn.Close()

		time.Sleep(50 * time.Millisecond)

		expected := []byte("hello")

		hub.Broadcast <- expected

		_, msg, err := conn.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, expected, msg)
	})

	t.Run("send to specific user", func(t *testing.T) {
		requireLocalListener(t)

		hub := notification.NewHub()
		go hub.Run()

		h, _ := NewNotificationHandler(hub)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
			h.ServeWS(w, r.WithContext(ctx))
		}))
		defer server.Close()

		conn, _, err := websocket.DefaultDialer.Dial(wsURL(server.URL), nil)
		assert.NoError(t, err)
		defer conn.Close()

		time.Sleep(50 * time.Millisecond)

		expected := []byte("private message")

		hub.SendToUser(userID, expected)

		_, msg, err := conn.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, expected, msg)
	})

	t.Run("client disconnect unregisters", func(t *testing.T) {
		requireLocalListener(t)

		hub := notification.NewHub()
		go hub.Run()

		h, _ := NewNotificationHandler(hub)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
			h.ServeWS(w, r.WithContext(ctx))
		}))
		defer server.Close()

		conn, _, err := websocket.DefaultDialer.Dial(wsURL(server.URL), nil)
		assert.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		conn.Close()

		time.Sleep(50 * time.Millisecond)

		hub.Mu.RLock()
		_, exists := hub.Clients[userID]
		hub.Mu.RUnlock()

		assert.False(t, exists)
	})
}
