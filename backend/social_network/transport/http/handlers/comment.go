package handlers

import (
	"encoding/json"
	"net/http"

	"social_network/internal/comment"
	"social_network/internal/notification"
	"social_network/internal/post"
	"social_network/internal/validation"
	"social_network/transport/http/dto"
	"social_network/transport/http/middleware"
	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CommentHandler struct {
	service             comment.CommentService
	postService         post.PostService
	notificationService notification.NotificationService
}

func NewCommentHandler(
	service comment.CommentService,
	postService post.PostService,
	notficationService notification.NotificationService,
) (*CommentHandler, error) {
	if service == nil {
		return nil, ErrCommentNil
	}
	if postService == nil {
		return nil, ErrPostNil
	}
	if notficationService == nil {
		return nil, ErrNotificationNil
	}

	return &CommentHandler{service: service, postService: postService, notificationService: notficationService}, nil
}

// @Summary      Создать комментарий
// @Description  Создает комментарий к посту и уведомляет автора поста
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        input body dto.CreateCommentRequest true "Comment data"
// @Success      201   {object} nil
// @Failure      400   {object} response.ErrorResponse
// @Failure      401   {object} response.ErrorResponse
// @Failure      500   {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /comments [post]
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	var request dto.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&request); err != nil {
		response.Error(w, err)
		return
	}

	err := h.service.Create(r.Context(), userID, request.PostID, request.Content)
	if err != nil {
		response.Error(w, err)
		return
	}

	post, err := h.postService.GetByID(r.Context(), request.PostID)
	if err != nil {
		response.Error(w, err)
		return
	}

	if post != nil && post.UserID != userID {
		go h.notificationService.NotifyComment(r.Context(), post.UserID, userID, request.PostID)
	}

	response.JSON(w, http.StatusCreated, nil)
}

// @Summary      Получить комментарий по ID
// @Description  Возвращает комментарий по его ID
// @Tags         comment
// @Produce      json
// @Param        id path string true "Comment ID"
// @Success      200   {object} dto.CommentResponse
// @Failure      400   {object} response.ErrorResponse
// @Failure      404   {object} response.ErrorResponse
// @Failure      500   {object} response.ErrorResponse
// @Router       /comments/{id} [get]
func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	comment, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, comment)
}

// @Summary      Список комментариев поста
// @Description  Возвращает список комментариев по ID поста
// @Tags         comment
// @Produce      json
// @Param        postID path string true "Post ID"
// @Success      200   {array} dto.CommentResponse
// @Failure      400   {object} response.ErrorResponse
// @Failure      500   {object} response.ErrorResponse
// @Router       /posts/{postID}/comments [get]
func (h *CommentHandler) ListByPostID(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	if postIDStr == "" {
		postIDStr = chi.URLParam(r, "id")
	}
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	comments, err := h.service.ListByPostID(r.Context(), postID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, comments)
}

// @Summary      Удалить комментарий
// @Description  Удаляет комментарий текущего пользователя по ID
// @Tags         comment
// @Produce      json
// @Param        id path string true "Comment ID"
// @Success      204   {object} nil
// @Failure      400   {object} response.ErrorResponse
// @Failure      401   {object} response.ErrorResponse
// @Failure      403   {object} response.ErrorResponse
// @Failure      500   {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /comments/{id} [delete]
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), userID, commentID); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
