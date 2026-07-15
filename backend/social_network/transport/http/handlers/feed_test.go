package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/post"
	"social_network/transport/http/middleware"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFeedHandler_GetFeed_Success(t *testing.T) {
	userID := uuid.New()

	mockService := new(MockFeedService)

	posts := []*post.Post{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Content:   "test post",
			CreatedAt: time.Now(),
		},
	}

	nextCursor := time.Now().Add(-time.Minute)

	mockService.
		On("GetFeed", mock.Anything, userID, mock.Anything, mock.Anything).
		Return(posts, nextCursor, nil)

	handler, _ := NewFeedHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

	w := httptest.NewRecorder()

	handler.GetFeed(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.NotNil(t, resp["posts"])
}

func TestFeedHandler_GetFeed_Unauthorized(t *testing.T) {
	mockService := new(MockFeedService)
	handler, _ := NewFeedHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	w := httptest.NewRecorder()

	handler.GetFeed(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFeedHandler_GetFeed_InvalidCursor(t *testing.T) {
	userID := uuid.New()

	mockService := new(MockFeedService)
	handler, _ := NewFeedHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/feed?cursor=invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

	w := httptest.NewRecorder()

	handler.GetFeed(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFeedHandler_GetFeed_ServiceError(t *testing.T) {
	userID := uuid.New()

	mockService := new(MockFeedService)

	mockService.
		On("GetFeed", mock.Anything, userID, mock.Anything, mock.Anything).
		Return([]*post.Post{}, time.Time{}, ErrServiceError)

	handler, _ := NewFeedHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

	w := httptest.NewRecorder()

	handler.GetFeed(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
