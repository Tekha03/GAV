package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gav/internal/follow"
	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
	"gav/internal/validation"
)

type FollowHandler struct {
	service follow.FollowService
}

func NewFollowHandler(service follow.FollowService) *FollowHandler {
	return &FollowHandler{service: service}
}

type followRequest struct {
	UserID uint	`json:"user_id" validate:"required"`
}

func (fh *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	authID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	var request followRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := validation.Validate(&request); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := fh.service.Follow(r.Context(), authID, request.UserID); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
