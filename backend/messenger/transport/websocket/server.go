package websocket

import (
	"encoding/json"
	"net/http"
	apperrors "shared/app_errors"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ChatHandler(hub *Hub, authVerifier func(w http.ResponseWriter, r *http.Request) (uuid.UUID, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 || parts[2] != "chat" {
			writeError(w, apperrors.New(apperrors.Validation, "invalid WebSocket path", apperrors.WithDetail("field", "path")))
			return
		}

		chatID, err := uuid.Parse(parts[3])
		if err != nil {
			writeError(w, apperrors.Wrap(apperrors.Validation, "invalid chat ID", err, apperrors.WithDetail("field", "chat_id")))
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
