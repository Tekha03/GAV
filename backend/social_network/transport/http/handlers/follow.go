package handlers

import (
	"net/http"

	"social_network/internal/follow"
	"social_network/internal/notification"
	"social_network/transport/http/middleware"
	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FollowHandler struct {
	service             follow.FollowService
	notificationService notification.NotificationService
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

// Follow godoc
// @Summary Follow user
// @Description Подписаться на пользователя
// @Tags follow
// @Accept json
// @Produce json
// @Param request body dto.FollowRequest true "User to follow"
// @Success 204 {object} nil
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /follow [post]
// @Security BearerAuth
func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	followerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	param := chi.URLParam(r, "userID")
	followingID, err := uuid.Parse(param)
	if err != nil {
		response.Error(w, err)
		return
	}

	newFollow, err := follow.NewFollow(followerID, followingID)
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Follow(r.Context(), *newFollow); err != nil {
		response.Error(w, err)
		return
	}

	go h.notificationService.NotifyFollow(r.Context(), followingID, followerID)

	response.JSON(w, http.StatusNoContent, nil)
}

// Unfollow godoc
// @Summary Unfollow user
// @Description Отписаться от пользователя
// @Tags follow
// @Produce json
// @Param userID path string true "User ID"
// @Success 204 {object} nil
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /follow/{userID} [delete]
// @Security BearerAuth
func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	followerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
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

// GetFollowers godoc
// @Summary Get followers
// @Description Получить список подписчиков
// @Tags follow
// @Produce json
// @Success 200 {array} follow.Follow
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /followers [get]
// @Security BearerAuth
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

// GetFollowing godoc
// @Summary Get following
// @Description Получить список подписок
// @Tags follow
// @Produce json
// @Success 200 {array} follow.Follow
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /following [get]
// @Security BearerAuth
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
