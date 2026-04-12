package notification

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHub_SendToUser(t *testing.T) {
	hub := NewHub()

	userID := uuid.New()

	client := &Client{
		UserID: userID,
		Send:   make(chan []byte, 1),
	}

	hub.Clients[userID] = client

	message := []byte("hello")

	hub.SendToUser(userID, message)

	select {
	case msg := <-client.Send:
		require.Equal(t, message, msg)
	default:
		t.Fatal("message not received")
	}
}
