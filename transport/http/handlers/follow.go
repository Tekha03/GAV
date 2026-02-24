package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/follow"
	"gav/internal/validation"
	"gav/transport/http/dto"
	"gav/transport/http/middleware"
	"gav/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FollowHandler struct {
	service follow.FollowService
}

func NewFollowHandler(service follow.FollowService) *FollowHandler {
	return &FollowHandler{service: service}
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	followerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	var following dto.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&following); err != nil {
		response.Error(w, err)
		return
	}

	if err := validation.Validate(&following); err != nil {
		response.Error(w, err)
		return
	}

	newFollow, err := follow.NewFollow(followerID, following.UserID)
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Follow(r.Context(), *newFollow); err != nil {
		response.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	followerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
	}

	param := chi.URLParam(r, "userID")
	followingID, err := uuid.Parse(param)
	if err != nil {
		response.Error(w, err)
		return
	}

	unfollow, err := follow.NewFollow(followerID, followingID)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	if err := h.service.Unfollow(r.Context(), *unfollow); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	followers, err := h.service.GetFollowers(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, followers)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	following, err := h.service.GetFollowing(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, following)
}
