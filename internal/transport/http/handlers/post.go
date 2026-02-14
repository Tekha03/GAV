package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gav/internal/post"
	"gav/internal/transport/http/dto"
	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
	"gav/internal/validation"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PostHandler struct {
	service post.PostService
}

func NewPostHandler(service post.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	var request dto.PostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, err)
		return
	}

	if err := validation.Validate(&request); err != nil {
		response.Error(w, err)
		return
	}

	post, err := h.service.Create(r.Context(), userID, request.Content)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.NewPostResponse(post))
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		response.Error(w, err)
		return
	}

	post, err := h.service.GetByID(r.Context(), uuid.UUID(id))
	if err != nil {
		response.Error(w, err)
		return
	}

	dtoPost := dto.NewPostResponse(post)
	response.JSON(w, http.StatusOK, dtoPost)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, errors.New("unauthorized"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	err = h.service.Delete(r.Context(), userID, uuid.UUID(id))
	if err != nil {
		response.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
