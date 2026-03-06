package handlers

import (
	"net/http"

	"gav/internal/media"
	"gav/internal/profile"
	"gav/transport/http/middleware"
	"gav/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UploadHandler struct {
	mediaService  media.MediaService
	profileService profile.ProfileService
}

func NewUploadHandler(mediaService media.MediaService) (*UploadHandler, error) {
	if mediaService == nil {
		return nil, ErrMediaNil
	}

	return &UploadHandler{mediaService: mediaService}, nil
}

func (h *UploadHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {	// 10MB max
		response.Error(w, err)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		response.Error(w, err)
		return
	}
	defer file.Close()

	url, err := h.mediaService.UploadImage(r.Context(), file, header, "avatars/"+userID.String())
	if err != nil {
		response.Error(w, err)
		return
	}

	profileID, err := uuid.Parse(chi.URLParam(r, "profileID"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}
	update := profile.UpdateProfileInput{ProfilePhotoUrl: &url}
	h.profileService.Update(r.Context(), profileID, update)

	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *UploadHandler) UploadPostImage(w http.ResponseWriter, r *http.Request) {
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

	url, err := h.mediaService.UploadImage(r.Context(), file, header, "posts/"+userID.String())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}


