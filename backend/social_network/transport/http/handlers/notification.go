package handlers

import (
	"net/http"

	"social_network/internal/notification"
	"social_network/transport/http/middleware"
	"social_network/transport/response"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool { return true },
}

type NotificationHandler struct {
	hub *notification.Hub
}

func NewNotificationHandler(hub *notification.Hub) (*NotificationHandler, error) {
	if hub == nil {
		return nil, notification.ErrEmptyHub
	}

	return &NotificationHandler{hub: hub}, nil
}

// ServeWS
// @Summary      WebSocket уведомления
// @Description  Устанавливает WebSocket соединение для получения уведомлений в реальном времени.
// @Description  После подключения клиент получает уведомления (лайки, комментарии, подписки).
// @Tags         notifications
// @Security     BearerAuth
// @Success      101  {string}  string  "WebSocket соединение установлено"
// @Failure      401  {object}  response.ErrorResponse  "Unauthorized"
// @Failure      500  {object}  response.ErrorResponse  "WebSocket upgrade error"
// @Router       /ws/notifications [get]
func (h *NotificationHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		response.Error(w, err)
		return
	}

	client := &notification.Client{
		UserID: userID,
		Conn:	conn,
		Send: 	make(chan []byte, 256),
	}

	h.hub.Register <- client

	go h.readPump(client)
	go h.writePump(client)
}

func (h *NotificationHandler) readPump(client *notification.Client) {
	defer func () {
		h.hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *NotificationHandler) writePump(client *notification.Client) {
	defer client.Conn.Close()
	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
