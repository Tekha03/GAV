package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/like"
	"social_network/internal/post"
	"social_network/transport/http/middleware"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLikeHandler_Add(t *testing.T) {
	userID := uuid.New()
	postID := uuid.New()

	t.Run("success with notification", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		input := like.Like{PostID: postID}

		mockLike.On("Add", mock.Anything, mock.Anything).Return(nil)

		mockPost.On("GetByID", mock.Anything, postID).
			Return(&post.Post{UserID: uuid.New()}, nil)

		mockNotif.
			On("NotifyLike", mock.Anything, mock.Anything, mock.Anything, postID).
			Maybe().
			Return(nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/likes", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Add(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("success without notification (self-like)", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		input := like.Like{PostID: postID}

		mockLike.On("Add", mock.Anything, mock.Anything).Return(nil)

		mockPost.On("GetByID", mock.Anything, postID).
			Return(&post.Post{UserID: userID}, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/likes", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Add(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		req := httptest.NewRequest(http.MethodPost, "/likes", nil)
		w := httptest.NewRecorder()

		h.Add(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		req := httptest.NewRequest(http.MethodPost, "/likes", bytes.NewBufferString("{bad"))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Add(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		input := like.Like{PostID: postID}

		mockLike.On("Add", mock.Anything, mock.Anything).Return(assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/likes", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Add(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("post service error", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		input := like.Like{PostID: postID}

		mockLike.On("Add", mock.Anything, mock.Anything).Return(nil)
		mockPost.On("GetByID", mock.Anything, postID).Return(nil, assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/likes", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Add(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})
}

func TestLikeHandler_Remove(t *testing.T) {
	userID := uuid.New()
	postID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		input := like.Like{PostID: postID}

		mockLike.On("Remove", mock.Anything, mock.Anything).Return(nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodDelete, "/likes", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Remove(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		input := like.Like{PostID: postID}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodDelete, "/likes", bytes.NewBuffer(body))

		w := httptest.NewRecorder()
		h.Remove(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		req := httptest.NewRequest(http.MethodDelete, "/likes", bytes.NewBufferString("{bad"))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Remove(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockLike := new(MockLikeService)
		mockPost := new(MockPostService)
		mockNotif := new(MockNotificationService)

		h, _ := NewLikeHandler(mockLike, mockPost, mockNotif)

		input := like.Like{PostID: postID}

		mockLike.On("Remove", mock.Anything, mock.Anything).Return(assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodDelete, "/likes", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Remove(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})
}
