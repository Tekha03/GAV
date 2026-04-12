package handlers

import (
	"encoding/json"
	"net/http"

	"social_network/internal/like"
	"social_network/internal/notification"
	"social_network/internal/post"
	"social_network/transport/http/middleware"
	"social_network/transport/response"
)

type LikeHandler struct {
	service 			like.LikeService
	postService			post.PostService
	notificationService  notification.NotificationService
}

func NewLikeHandler(
	service like.LikeService,
	postService post.PostService,
	notificationService notification.NotificationService,
) (*LikeHandler, error) {
	if service == nil {
		return nil, ErrLikeNil
	}
	if postService == nil {
		return nil, ErrPostNil
	}
	if notificationService == nil {
		return nil, ErrNotificationNil
	}

	return &LikeHandler{service: service, postService: postService, notificationService: notificationService}, nil
}

// Add
// @Summary      Добавить лайк
// @Description  Добавляет лайк к посту текущего пользователя. Если лайк на чужой пост — создается уведомление.
// @Tags         likes
// @Accept       json
// @Produce      json
// @Param        like  body      like.Like  true  "Данные лайка (PostID)"
// @Success      204   {object}  nil        "Лайк успешно добавлен"
// @Failure      400   {object}  response.ErrorResponse  "Некорректный ввод"
// @Failure      401   {object}  response.ErrorResponse  "Неавторизованный пользователь"
// @Failure      500   {object}  response.ErrorResponse  "Внутренняя ошибка сервиса"
// @Router       /likes [post]
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

	post, err := h.postService.GetByID(r.Context(), like.PostID)
	if err != nil {
		response.Error(w, err)
		return
	}

	if post != nil && post.UserID != userID {
		go h.notificationService.NotifyLike(r.Context(), post.UserID, userID, like.PostID)
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// Remove
// @Summary      Удалить лайк
// @Description  Удаляет лайк текущего пользователя с поста.
// @Tags         likes
// @Accept       json
// @Produce      json
// @Param        like  body      like.Like  true  "Данные лайка (PostID)"
// @Success      204   {object}  nil        "Лайк успешно удален"
// @Failure      400   {object}  response.ErrorResponse  "Некорректный ввод"
// @Failure      401   {object}  response.ErrorResponse  "Неавторизованный пользователь"
// @Failure      500   {object}  response.ErrorResponse  "Внутренняя ошибка сервиса"
// @Router       /likes [delete]
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

	response.JSON(w, http.StatusNoContent, nil)
}
