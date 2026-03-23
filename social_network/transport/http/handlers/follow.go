package handlers

import (
	"encoding/json"
	"net/http"

	"social_network/internal/follow"
	"social_network/internal/notification"
	"social_network/internal/validation"
	"social_network/transport/http/dto"
	"social_network/transport/http/middleware"
	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FollowHandler struct {
	service 			follow.FollowService
	notificationService  notification.NotificationService
}

func NewFollowHandler(service follow.FollowService, notificationService notification.NotificationService) (*FollowHandler, error) {
	if service == nil {
		return nil, ErrFollowNil
	}
	if notificationService == nil {
		return nil, ErrNotificationNil
	}

	return &FollowHandler{service: service, notificationService: notificationService}, nil
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

	go h.notificationService.NotifyFollow(r.Context(), following.UserID, followerID)

	response.JSON(w, http.StatusNoContent, nil)
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