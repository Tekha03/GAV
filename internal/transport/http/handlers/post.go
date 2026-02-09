package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gav/internal/post"
	"gav/internal/transport/http/dto"
	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
	"gav/internal/validation"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	service post.Service
}

func NewPostHandler(service post.Service) *PostHandler {
	return &PostHandler{service: service}
}

func (ph *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, errors.New("unauthorized"))
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

	post, err := ph.service.Create(r.Context(), userID, request.Content)
	if err != nil {
		response.Error(w, err)
		return
	}

	dtoPost := dto.NewPostResponse(post)
	response.JSON(w, http.StatusCreated, dtoPost)
}

func (ph *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		response.Error(w, err)
		return
	}

	post, err := ph.service.GetByID(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}

	dtoPost := dto.NewPostResponse(post)
	response.JSON(w, http.StatusOK, dtoPost)
}

func (ph *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, errors.New("unauthorized"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, err)
		return
	}

	err = ph.service.Delete(r.Context(), userID, uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
