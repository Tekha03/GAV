package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"messenger/internal/model"
	"messenger/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type chatDTO struct {
	ID        uuid.UUID `json:"id"`
	IsGroup   bool      `json:"is_group"`
	Title     string    `json:"title"`
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}

type messageDTO struct {
	ID          uuid.UUID       `json:"id"`
	ChatID      uuid.UUID       `json:"chat_id"`
	SenderID    uuid.UUID       `json:"sender_id"`
	Text        *string         `json:"text"`
	ReplyToID   *uuid.UUID      `json:"reply_to_id"`
	CreatedAt   time.Time       `json:"created_at"`
	EditedAt    *time.Time      `json:"edited_at"`
	Attachments []attachmentDTO `json:"attachments"`
}

type chatMemberDTO struct {
	ChatID            uuid.UUID `json:"chat_id"`
	UserID            uuid.UUID `json:"user_id"`
	Role              string    `json:"role"`
	Muted             bool      `json:"muted"`
	LastReadMessageID uuid.UUID `json:"last_read_message_id,omitempty"`
}

type attachmentDTO struct {
	ID        uuid.UUID `json:"id,omitempty"`
	MessageID uuid.UUID `json:"message_id,omitempty"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
}

type createPrivateChatRequest struct {
	UserID1 uuid.UUID `json:"user_id_1"`
	UserID2 uuid.UUID `json:"user_id_2"`
}

type createGroupChatRequest struct {
	Title     string      `json:"title"`
	CreatorID uuid.UUID   `json:"creator_id"`
	MemberIDs []uuid.UUID `json:"member_ids"`
}

type sendMessageRequest struct {
	SenderID    uuid.UUID       `json:"sender_id"`
	Text        *string         `json:"text"`
	ReplyToID   *uuid.UUID      `json:"reply_to_id"`
	Attachments []attachmentDTO `json:"attachments"`
}

func NewHTTPServer(addr string, chatService service.Service) *http.Server {
	r := chi.NewRouter()
	r.Use(corsMiddleware())

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/chats/private", createPrivateChat(chatService))
		r.Post("/chats/group", createGroupChat(chatService))
		r.Get("/users/{user_id}/chats", getUserChats(chatService))
		r.Get("/chats/{chat_id}", getChat(chatService))
		r.Get("/chats/{chat_id}/members", getChatMembers(chatService))
		r.Get("/chats/{chat_id}/messages", getMessages(chatService))
		r.Post("/chats/{chat_id}/messages", sendMessage(chatService))
		r.Post("/chats/{chat_id}/read", markAsRead(chatService))
	})

	return &http.Server{Addr: addr, Handler: r}
}

func getChatMembers(chatService service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, err := parsePathUUID(r, "chat_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		members, err := chatService.GetChatMembers(r.Context(), chatID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		items := make([]chatMemberDTO, 0, len(members))
		for _, member := range members {
			items = append(items, toChatMemberDTO(member))
		}

		writeJSON(w, http.StatusOK, map[string][]chatMemberDTO{"members": items})
	}
}

func createPrivateChat(chatService service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPrivateChatRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		chat, err := chatService.CreatePrivateChat(r.Context(), req.UserID1, req.UserID2)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusCreated, map[string]chatDTO{"chat": toChatDTO(chat)})
	}
}

func createGroupChat(chatService service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createGroupChatRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		chat, err := chatService.CreateGroupChat(r.Context(), req.Title, req.CreatorID, req.MemberIDs)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusCreated, map[string]chatDTO{"chat": toChatDTO(chat)})
	}
}

func getUserChats(chatService service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := parsePathUUID(r, "user_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		chats, err := chatService.GetUserChats(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		items := make([]chatDTO, 0, len(chats))
		for _, chat := range chats {
			items = append(items, toChatDTO(chat))
		}

		writeJSON(w, http.StatusOK, map[string][]chatDTO{"chats": items})
	}
}

func getChat(chatService service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, err := parsePathUUID(r, "chat_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		chat, err := chatService.GetChatByID(r.Context(), chatID)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]chatDTO{"chat": toChatDTO(chat)})
	}
}

func getMessages(chatService service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, err := parsePathUUID(r, "chat_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		limit := 50
		if raw := r.URL.Query().Get("limit"); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		var cursor *uuid.UUID
		if raw := r.URL.Query().Get("before"); raw != "" {
			parsed, err := uuid.Parse(raw)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			cursor = &parsed
		}

		messages, err := chatService.GetMessages(r.Context(), chatID, limit, cursor)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		items := make([]messageDTO, 0, len(messages))
		for _, message := range messages {
			items = append(items, toMessageDTO(message))
		}

		writeJSON(w, http.StatusOK, map[string][]messageDTO{"messages": items})
	}
}

func sendMessage(chatService service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, err := parsePathUUID(r, "chat_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		var req sendMessageRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		attachments := make([]model.AttachmentInput, 0, len(req.Attachments))
		for _, attachment := range req.Attachments {
			attachments = append(attachments, model.AttachmentInput{
				URL:      attachment.URL,
				Type:     model.AttachmentType(attachment.Type),
				FileName: attachment.FileName,
				FileSize: attachment.FileSize,
			})
		}

		message, err := chatService.SendMessage(r.Context(), model.SendMessageInput{
			ChatID:      chatID,
			SenderID:    req.SenderID,
			Text:        req.Text,
			ReplyToID:   req.ReplyToID,
			Attachments: attachments,
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusCreated, map[string]messageDTO{
			"message": toMessageDTO(message),
		})
	}
}

func markAsRead(chatService service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID, err := parsePathUUID(r, "chat_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		type request struct {
			UserID uuid.UUID `json:"user_id"`
		}
		var req request
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		if err := chatService.MarkAsRead(r.Context(), chatID, req.UserID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func parsePathUUID(r *http.Request, key string) (uuid.UUID, error) {
	raw := chi.URLParam(r, key)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", key)
	}
	return id, nil
}

func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func toChatDTO(chat *model.Chat) chatDTO {
	return chatDTO{
		ID:        chat.ID,
		IsGroup:   chat.IsGroup,
		Title:     chat.Title,
		PhotoURL:  chat.PhotoURL,
		CreatedAt: chat.CreatedAt,
	}
}

func toMessageDTO(message *model.Message) messageDTO {
	attachments := make([]attachmentDTO, 0, len(message.Attachments))
	for _, attachment := range message.Attachments {
		attachments = append(attachments, attachmentDTO{
			ID:        attachment.ID,
			MessageID: attachment.MessageID,
			URL:       attachment.URL,
			Type:      string(attachment.Type),
			FileName:  attachment.FileName,
			FileSize:  attachment.FileSize,
		})
	}

	return messageDTO{
		ID:          message.ID,
		ChatID:      message.ChatID,
		SenderID:    message.SenderID,
		Text:        message.Text,
		ReplyToID:   message.ReplyToID,
		CreatedAt:   message.CreatedAt,
		EditedAt:    message.EditedAt,
		Attachments: attachments,
	}
}

func toChatMemberDTO(member *model.ChatMember) chatMemberDTO {
	return chatMemberDTO{
		ChatID:            member.ChatID,
		UserID:            member.UserID,
		Role:              string(member.Role),
		Muted:             member.Muted,
		LastReadMessageID: member.LastReadMessageID,
	}
}

func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
