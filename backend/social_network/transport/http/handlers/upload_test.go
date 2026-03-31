package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"social_network/internal/profile"
	"social_network/transport/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUploadHandler(t *testing.T) *UploadHandler {
	mediaMock := new(MockUploadService)
	profileMock := new(MockProfileService)

	h, err := NewUploadHandler(mediaMock)
	if err != nil {
		t.Fatal(err)
	}
	h.ProfileService = profileMock
	return h
}

func TestUploadHandler_UploadAvatar_Success(t *testing.T) {
	handler := setupUploadHandler(t)
	userID := uuid.New()
	profileID := uuid.New()
	expectedURL := "https://cdn.example.com/avatar.jpg"

	mediaMock := handler.MediaService.(*MockUploadService)
	mediaMock.On("UploadImage", mock.Anything, mock.Anything, mock.Anything, "avatars/"+userID.String()).
		Return(expectedURL, nil)

	profileMock := handler.ProfileService.(*MockProfileService)
	profileMock.On("Update", mock.Anything, profileID,
		profile.UpdateProfileInput{ProfilePhotoUrl: &expectedURL}).Return(nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	part.Write([]byte("dummy image content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("profileID", profileID.String())
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)

	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.UploadAvatar(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), expectedURL)

	mediaMock.AssertExpectations(t)
	profileMock.AssertExpectations(t)
}

func TestUploadHandler_UploadAvatar_Unauthorized(t *testing.T) {
	handler := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	handler.UploadAvatar(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUploadHandler_UploadPostImage_Success(t *testing.T) {
	handler := setupUploadHandler(t)

	userID := uuid.New()
	expectedURL := "https://cdn.example.com/post.jpg"

	mediaMock := handler.MediaService.(*MockUploadService)
	mediaMock.On("UploadImage", mock.Anything, mock.Anything, mock.Anything, "posts/"+userID.String()).
		Return(expectedURL, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "image.jpg")
	part.Write([]byte("dummy image content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/post-image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

	w := httptest.NewRecorder()
	handler.UploadPostImage(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), expectedURL)
}

func TestUploadHandler_UploadPostImage_Unauthorized(t *testing.T) {
	handler := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/post-image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	handler.UploadPostImage(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
