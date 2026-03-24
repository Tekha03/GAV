package handlers

import (
	"encoding/json"
	"net/http"

	"social_network/internal/media"
	"social_network/internal/post"
	"social_network/internal/validation"
	"social_network/transport/http/dto"
	"social_network/transport/http/middleware"
	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PostHandler struct {
	service 	 		post.PostService
	mediaService 		media.MediaService
}

func NewPostHandler(service post.PostService, mediaService media.MediaService) (*PostHandler, error) {
	if service == nil {
		return nil, ErrPostNil
	}
	if mediaService == nil {
		return nil, ErrMediaNil
	}

	return &PostHandler{service: service, mediaService: mediaService}, nil
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, err)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		response.Error(w, err)
		return
	}
	defer file.Close()

	imageUrl, err := h.mediaService.UploadImage(r.Context(), file, header, "posts/"+userID.String())
	if err != nil {
		response.Error(w, err)
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

	post, err := h.service.Create(r.Context(), userID, request.Content, imageUrl)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.PostResponse{
		ID: post.ID,
		AuthorID: post.UserID,
		Content: post.Content,
		ImageUrl: post.ImageUrl,
		CreatedAt: post.CreatedAt,
	})
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

	response.JSON(w, http.StatusOK, dto.PostResponse{
		ID: post.ID,
		AuthorID:
		post.UserID,
		Content: post.Content,
		CreatedAt: post.CreatedAt,
	})
}

func (h *PostHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	list, err := h.service.ListByUser(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
	}

	response.JSON(w, http.StatusOK, list)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
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