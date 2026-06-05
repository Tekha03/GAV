package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/follow"
	"social_network/transport/http/dto"
	"social_network/transport/http/middleware"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFollowHandler_Follow(t *testing.T) {
	followerID := uuid.New()
	followingID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		input := dto.FollowRequest{UserID: followingID}

		mockSvc.On("Follow", mock.Anything, mock.Anything).Return(nil)
		mockNotif.On("NotifyFollow", mock.Anything, followingID, followerID).Return()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, followerID))

		w := httptest.NewRecorder()
		h.Follow(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		req := httptest.NewRequest(http.MethodPost, "/follow", nil)
		w := httptest.NewRecorder()

		h.Follow(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		req := httptest.NewRequest(http.MethodPost, "/follow", bytes.NewBufferString("{bad"))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, followerID))

		w := httptest.NewRecorder()
		h.Follow(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		input := dto.FollowRequest{UserID: followingID}

		mockSvc.On("Follow", mock.Anything, mock.Anything).Return(assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, followerID))

		w := httptest.NewRecorder()
		h.Follow(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})
}

func TestFollowHandler_Unfollow(t *testing.T) {
	followerID := uuid.New()
	followingID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		mockSvc.On("Unfollow", mock.Anything, mock.Anything).Return(nil)

		r := chi.NewRouter()
		r.Delete("/follow/{userID}", h.Unfollow)

		req := httptest.NewRequest(http.MethodDelete, "/follow/"+followingID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, followerID))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		r := chi.NewRouter()
		r.Delete("/follow/{userID}", h.Unfollow)

		req := httptest.NewRequest(http.MethodDelete, "/follow/invalid", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, followerID))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})
}

func TestFollowHandler_GetFollowers(t *testing.T) {
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		mockSvc.On("GetFollowers", mock.Anything, userID).
			Return([]follow.Follow{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/followers", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.GetFollowers(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		req := httptest.NewRequest(http.MethodGet, "/followers", nil)
		w := httptest.NewRecorder()

		h.GetFollowers(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestFollowHandler_GetFollowing(t *testing.T) {
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		mockSvc.On("GetFollowing", mock.Anything, userID).
			Return([]follow.Follow{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/following", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.GetFollowing(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(MockFollowService)
		mockNotif := new(MockNotificationService)

		h, _ := NewFollowHandler(mockSvc, mockNotif)

		req := httptest.NewRequest(http.MethodGet, "/following", nil)
		w := httptest.NewRecorder()

		h.GetFollowing(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}
