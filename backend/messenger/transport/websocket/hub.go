package websocket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type OutgoingMessage struct {
	ChatID uuid.UUID
	Data   []byte
}

type Client struct {
	UserID uuid.UUID
	ChatID uuid.UUID
	Conn   *websocket.Conn
	Send   chan OutgoingMessage
}

type Hub struct {
	clients    map[uuid.UUID][]*Client
	Register   chan *Client
	Unregister chan *Client
	sendToChat chan OutgoingMessage
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		sendToChat: make(chan OutgoingMessage, 1024),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client.ChatID] = append(h.clients[client.ChatID], client)
		case client := <-h.Unregister:
			h.removeClient(client)
		case message := <-h.sendToChat:
			for _, client := range h.clients[message.ChatID] {
				select {
				case client.Send <- message:
				default:
					h.removeClient(client)
				}
			}
		}
	}
}

func (h *Hub) SendToChat(chatID uuid.UUID, data []byte) {
	h.sendToChat <- OutgoingMessage{ChatID: chatID, Data: data}
}

func (h *Hub) removeClient(client *Client) {
	clients := h.clients[client.ChatID]
	for i, existingClient := range clients {
		if existingClient == client {
			h.clients[client.ChatID] = append(clients[:i], clients[i+1:]...)
			close(client.Send)
			return
		}
	}
}
