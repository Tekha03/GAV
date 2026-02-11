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
	service follow.Service
}

func NewFollowHandler(service follow.Service) *FollowHandler {
	return &FollowHandler{service: service}
}

type followRequest struct {
	UserID uint	`json:"user_id" validate:"required"`
}

func (fh *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	authID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, errors.New("unauthorized"))
		return
	}

	var request followRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, err)
		return
	}

	if err := validation.Validate(&request); err != nil {
		response.Error(w, err)
		return
	}

	newFollow := follow.NewFollow(authID, request.UserID)
	if err := fh.service.Follow(r.Context(), *newFollow); err != nil {
		response.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
