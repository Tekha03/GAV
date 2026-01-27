package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"gav/internal/post"
	"gav/internal/transport/httpserver/middleware"
	"gav/internal/transport/response"
)

type PostHandler struct {
	service post.Service
}

func (ph *PostHandler) NewPostHandler(service post.Service) *PostHandler {
	return &PostHandler{service: service}
}

type createPostRequest struct {
	PostID	uint	`json:"post_id"`
	Content	string	`json:"content"`
}

func (ph *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	var request createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post, err := ph.service.Create(r.Context(), userID, request.PostID, request.Content)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.JSON(w, http.StatusCreated, post)
}

func (ph *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post, err := ph.service.GetByID(r.Context(), uint(id))
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

func (ph *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	err = ph.service.Delete(r.Context(), userID, uint(id))
	if err != nil {
		response.Error(w, http.StatusForbidden, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
