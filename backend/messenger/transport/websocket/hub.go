package websocket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID uuid.UUID
	ChatID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients    map[uuid.UUID][]*Client
	Register   chan *Client
	Unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client.ChatID] = append(h.clients[client.ChatID], client)
		case client := <-h.Unregister:
			clients := h.clients[client.ChatID]
			for i, c := range clients {
				if c == client {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			h.clients[client.ChatID] = clients
			close(client.Send)
		case message := <-h.broadcast:
			for _, clients := range h.clients {
				for _, client := range clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						h.Unregister <- client
					}
				}
			}
		}
	}
}

func (h *Hub) SendToChat(chatID uuid.UUID, data []byte) {
	clients := h.clients[chatID]
	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			h.Unregister <- client
		}
	}
}
