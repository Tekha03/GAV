package websocket

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ChatHandler(hub *Hub, authVerifier func(w http.ResponseWriter, r *http.Request) (uuid.UUID, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 || parts[2] != "chat" {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}

		chatID, err := uuid.Parse(parts[3])
		if err != nil {
			http.Error(w, "bad chat_id", http.StatusBadRequest)
			return
		}

		userID, err := authVerifier(w, r)
		if err != nil {
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := &Client{
			UserID: userID,
			ChatID: chatID,
			Conn:   conn,
			Send:   make(chan []byte),
		}

		hub.Register <- client

		go func() {
			for {
				_, _, err := client.Conn.ReadMessage()
				if err != nil {
					hub.Unregister <- client
					break
				}
			}
		}()

		go func() {
			for {
				select {
				case message, ok := <-client.Send:
					if !ok {
						return
					}
					err := client.Conn.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						hub.Unregister <- client
						return
					}
				}
			}
		}()
	}
}
