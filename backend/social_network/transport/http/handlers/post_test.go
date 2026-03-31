package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"social_network/internal/post"
	"social_network/transport/http/middleware"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type postTestEnv struct {
	handler      *PostHandler
	postService  *MockPostService
	mediaService *MockMediaService
}

func setupPostHandler(t *testing.T) *postTestEnv {
	postService := &MockPostService{}
	mediaService := &MockMediaService{}

	handler, err := NewPostHandler(postService, mediaService)
	require.NoError(t, err)

	return &postTestEnv{
		handler:      handler,
		postService:  postService,
		mediaService: mediaService,
	}
}

func TestPostHandler_Create_Success(t *testing.T) {
	env := setupPostHandler(t)

	userID := uuid.New()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("file-content"))

	writer.WriteField("content", "hello")

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/posts", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	env.mediaService.
		On("UploadImage", mock.Anything, mock.Anything, mock.Anything, "posts/"+userID.String()).
		Return("image-url", nil)

	env.postService.
		On("Create", mock.Anything, userID, mock.Anything, "image-url").
		Return(&post.Post{
			ID:        uuid.New(),
			UserID:    userID,
			Content:   "hello",
			ImageUrl:  "image-url",
			CreatedAt: time.Now(),
		}, nil)

	env.handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPostHandler_Create_Unauthorized(t *testing.T) {
	env := setupPostHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/posts", nil)
	w := httptest.NewRecorder()

	env.handler.Create(w, req)

	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestPostHandler_Create_UploadError(t *testing.T) {
	env := setupPostHandler(t)

	userID := uuid.New()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("file"))

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/posts", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	env.mediaService.
		On("UploadImage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", assert.AnError)

	env.handler.Create(w, req)

	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestPostHandler_GetByID_Success(t *testing.T) {
	env := setupPostHandler(t)

	postID := uuid.New()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", postID.String())

	req := httptest.NewRequest(http.MethodGet, "/posts/"+postID.String(), nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	env.postService.
		On("GetByID", mock.Anything, postID).
		Return(&post.Post{
			ID:        postID,
			UserID:    uuid.New(),
			Content:   "text",
			CreatedAt: time.Now(),
		}, nil)

	env.handler.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_GetByID_InvalidID(t *testing.T) {
	env := setupPostHandler(t)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")

	req := httptest.NewRequest(http.MethodGet, "/posts/invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	env.handler.GetByID(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestPostHandler_ListByUser_Success(t *testing.T) {
	env := setupPostHandler(t)

	userID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	env.postService.
		On("ListByUser", mock.Anything, userID).
		Return([]*post.Post{}, nil)

	env.handler.ListByUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_Delete_Unauthorized(t *testing.T) {
	env := setupPostHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	w := httptest.NewRecorder()

	env.handler.Delete(w, req)

	assert.NotEqual(t, http.StatusNoContent, w.Code)
}

func TestPostHandler_Delete_Success(t *testing.T) {
	env := setupPostHandler(t)

	userID := uuid.New()
	postID := uuid.New()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", postID.String())

	req := httptest.NewRequest(http.MethodDelete, "/posts/"+postID.String(), nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	env.postService.
		On("Delete", mock.Anything, userID, postID).
		Return(nil)

	env.handler.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
