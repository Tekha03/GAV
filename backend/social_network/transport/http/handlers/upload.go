package handlers

import (
	"net/http"

	"social_network/internal/media"
	"social_network/internal/profile"
	"social_network/transport/http/middleware"
	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UploadHandler struct {
	MediaService  media.MediaService
	ProfileService profile.ProfileService
}

func NewUploadHandler(mediaService media.MediaService) (*UploadHandler, error) {
	if mediaService == nil {
		return nil, ErrMediaNil
	}

	return &UploadHandler{MediaService: mediaService}, nil
}

// @Summary      Загрузка аватара пользователя
// @Description  Загружает фото аватара (jpg/png/webp, max 5MB)
// @Tags         upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        avatar  formData  file  true  "Фото аватара"
// @Security     BearerAuth
// @Success      200   {object} map[string]string
// @Failure      400   {object} response.ErrorResponse
// @Failure      401   {object} response.ErrorResponse
// @Router       /upload/avatar [post]
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

	url, err := h.MediaService.UploadImage(r.Context(), file, header, "avatars/"+userID.String())
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
	h.ProfileService.Update(r.Context(), profileID, update)

	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}

// @Summary      Загрузка изображения для поста
// @Description  Загружает изображение для поста. Допустимые форматы: jpg, png, webp. Максимальный размер 5MB.
// @Tags         upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file  true  "Изображение для поста"
// @Security     BearerAuth
// @Success      200  {object} map[string]string  "Ссылка на загруженное изображение"
// @Failure      400  {object} response.ErrorResponse  "Неверный запрос"
// @Failure      401  {object} response.ErrorResponse  "Неавторизованный"
// @Router       /upload/post-image [post]
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

	url, err := h.MediaService.UploadImage(r.Context(), file, header, "posts/"+userID.String())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}
