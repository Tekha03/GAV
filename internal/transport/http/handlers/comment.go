package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/comment"
	"gav/internal/transport/http/dto"
	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
	"gav/internal/validation"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CommentHandler struct {
	service comment.CommentService
}

func NewCommentHandler(service comment.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

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

	response.JSON(w, http.StatusCreated, nil)
}

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

func (h *CommentHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	comments, err := h.service.GetByPostID(r.Context(), postID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, comments)
}

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
