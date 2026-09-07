package websocket

import (
	"encoding/json"
	"messenger/internal/service"
	"net/http"
	apperrors "shared/app_errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ChatHandler(hub *Hub, chatService service.Service, authVerifier func(r *http.Request) (uuid.UUID, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, err := uuid.Parse(chi.URLParam(r, "chat_id"))
		if err != nil {
			writeError(w, apperrors.Wrap(apperrors.Validation, "invalid chat ID", err, apperrors.WithDetail("field", "chat_id")))
			return
		}

		userID, err := authVerifier(r)
		if err != nil {
			writeError(w, err)
			return
		}

		if err := chatService.RequireChatMember(r.Context(), chatID, userID); err != nil {
			writeError(w, err)
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
			Send:   make(chan OutgoingMessage, 256),
		}

		hub.Register <- client

		go client.readPump(hub)
		go client.writePump(hub)
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (c *Client) writePump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message.Data); err != nil {
			return
		}
	}
}

func writeError(w http.ResponseWriter, err error) {
	appErr := apperrors.Normalize(err)
	message := appErr.Message
	if appErr.Category == apperrors.CategoryInternal {
		message = "internal server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus(appErr.Category))
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code": appErr.Code, "category": appErr.Category,
			"message": message, "details": appErr.Details,
		},
	})
}

func httpStatus(category apperrors.Category) int {
	switch category {
	case apperrors.CategoryValidation:
		return http.StatusBadRequest
	case apperrors.CategoryUnauthenticated:
		return http.StatusUnauthorized
	case apperrors.CategoryPermissionDenied:
		return http.StatusForbidden
	case apperrors.CategoryNotFound:
		return http.StatusNotFound
	case apperrors.CategoryConflict:
		return http.StatusConflict
	case apperrors.CategoryUnavailable:
		return http.StatusServiceUnavailable
	case apperrors.CategoryUnsupported:
		return http.StatusNotImplemented
	case apperrors.CategoryCancelled:
		return 499
	default:
		return http.StatusInternalServerError
	}
}
