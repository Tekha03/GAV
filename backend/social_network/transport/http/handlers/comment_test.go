package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/comment"
	"social_network/internal/post"
	"social_network/transport/http/dto"
	"social_network/transport/http/middleware"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentHandler_Create(t *testing.T) {
	userID := uuid.New()
	postID := uuid.New()
	dtoReq := dto.CreateCommentRequest{PostID: postID, Content: "hello"}

	t.Run("success", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)

		mockComment.On("Create", mock.Anything, userID, postID, "hello").Return(nil)
		mockPost.On("GetByID", mock.Anything, postID).Return(&post.Post{ID: postID}, nil)
		mockNotif.On("NotifyComment", mock.Anything, mock.Anything, userID, postID).Return(nil)

		body, _ := json.Marshal(dtoReq)
		req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)

		req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBufferString("{bad"))
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.Create(w, req)

		assert.NotEqual(t, http.StatusCreated, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)

		mockComment.On("Create", mock.Anything, userID, postID, "hello").Return(ErrFail)

		body, _ := json.Marshal(dtoReq)
		req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.Create(w, req)

		assert.NotEqual(t, http.StatusCreated, w.Code)
	})
}

func TestCommentHandler_GetByID(t *testing.T) {
	commentID := uuid.New()
	commentResp := &comment.Comment{ID: commentID, Content: "hi"}

	t.Run("success", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)
		mockComment.On("GetByID", mock.Anything, commentID).Return(commentResp, nil)

		r := chi.NewRouter()
		r.Get("/comments/{id}", h.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/comments/"+commentID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "hi")
	})

	t.Run("invalid id", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)

		r := chi.NewRouter()
		r.Get("/comments/{id}", h.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/comments/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestCommentHandler_ListByPostID(t *testing.T) {
	postID := uuid.New()
	comments := []comment.Comment{{ID: uuid.New(), Content: "c1"}}

	t.Run("success", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)
		mockComment.On("ListByPostID", mock.Anything, postID).Return(comments, nil)

		r := chi.NewRouter()
		r.Get("/posts/{postID}/comments", h.ListByPostID)

		req := httptest.NewRequest(http.MethodGet, "/posts/"+postID.String()+"/comments", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "c1")
	})

	t.Run("invalid post id", func(t *testing.T) {
		mockComment := new(MockCommentService)
		h, _ := NewCommentHandler(mockComment, nil, nil)

		r := chi.NewRouter()
		r.Get("/posts/{postID}/comments", h.ListByPostID)

		req := httptest.NewRequest(http.MethodGet, "/posts/invalid/comments", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestCommentHandler_Delete(t *testing.T) {
	userID := uuid.New()
	commentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)
		mockComment.On("Delete", mock.Anything, userID, commentID).Return(nil)

		r := chi.NewRouter()
		r.Delete("/comments/{id}", h.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/comments/"+commentID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)

		r := chi.NewRouter()
		r.Delete("/comments/{id}", h.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/comments/invalid", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockComment := new(MockCommentService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)
		h, _ := NewCommentHandler(mockComment, mockPost, mockNotif)
		mockComment.On("Delete", mock.Anything, userID, commentID).Return(ErrFail)

		r := chi.NewRouter()
		r.Delete("/comments/{id}", h.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/comments/"+commentID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})
}
