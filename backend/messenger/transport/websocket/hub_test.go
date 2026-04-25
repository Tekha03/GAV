package websocket

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHubRegisterBroadcastSendAndUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	chatID := uuid.New()
	first := &Client{UserID: uuid.New(), ChatID: chatID, Send: make(chan []byte, 2)}
	second := &Client{UserID: uuid.New(), ChatID: uuid.New(), Send: make(chan []byte, 2)}

	hub.Register <- first
	hub.Register <- second
	time.Sleep(20 * time.Millisecond)

	hub.broadcast <- []byte("broadcast")
	select {
	case message := <-first.Send:
		if string(message) != "broadcast" {
			t.Fatalf("unexpected first broadcast message: %s", string(message))
		}
	case <-time.After(time.Second):
		t.Fatal("first client did not receive broadcast")
	}

	select {
	case message := <-second.Send:
		if string(message) != "broadcast" {
			t.Fatalf("unexpected second broadcast message: %s", string(message))
		}
	case <-time.After(time.Second):
		t.Fatal("second client did not receive broadcast")
	}

	hub.SendToChat(chatID, []byte("chat"))
	select {
	case message := <-first.Send:
		if string(message) != "chat" {
			t.Fatalf("unexpected chat message: %s", string(message))
		}
	case <-time.After(time.Second):
		t.Fatal("first client did not receive chat message")
	}

	select {
	case <-second.Send:
		t.Fatal("second client should not receive chat-specific message")
	case <-time.After(100 * time.Millisecond):
	}

	hub.Unregister <- first
	select {
	case _, ok := <-first.Send:
		if ok {
			t.Fatal("first send channel should be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("first client was not unregistered")
	}
}
