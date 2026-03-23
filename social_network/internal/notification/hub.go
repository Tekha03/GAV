package notification

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID	uuid.UUID
	Conn	*websocket.Conn
	Send 	chan []byte
}

type Hub struct {
	clients		map[uuid.UUID]*Client
	Register	chan *Client
	Unregister	chan *Client
	broadcast	chan []byte
	mu			sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: 	make(map[uuid.UUID]*Client),
		Register: 	make(chan *Client),
		Unregister: make(chan *Client),
		broadcast: 	make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if c, ok := h.clients[client.UserID]; ok && c == client {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, message []byte) {
	h.mu.Lock()
	if client, ok := h.clients[userID]; ok {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, userID)
		}
	}
	h.mu.Unlock()
}