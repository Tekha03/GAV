package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/like"
	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
)

type LikeHandler struct {
	service like.LikeService
}

func NewLikeHandler(service like.LikeService) *LikeHandler {
	return &LikeHandler{service: service}
}

func (h *LikeHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	var like like.Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	like.UserID = userID

	if err := h.service.Add(r.Context(), like); err != nil {
		response.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LikeHandler) Remove(w http.ResponseWriter, r *http.Request) {
	var like like.Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	like.UserID = userID

	if err := h.service.Remove(r.Context(), like); err != nil {
		response.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
