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
)

func TestFeedHandler_GetFeed_Success(t *testing.T) {
	userID := uuid.New()

	mockService := &MockFeedService{
		GetFeedFn: func(ctx context.Context, uid uuid.UUID, before time.Time, limit int) ([]*post.Post, time.Time, error) {
			return []*post.Post{
				{
					ID:        uuid.New(),
					UserID:    uid,
					Content:   "test post",
					CreatedAt: time.Now(),
				},
			}, time.Now().Add(-time.Minute), nil
		},
	}

	handler, _ := NewFeedHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/feed", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.GetFeed(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp["posts"] == nil {
		t.Fatal("expected posts in response")
	}
}

func TestFeedHandler_GetFeed_Unauthorized(t *testing.T) {
	handler, _ := NewFeedHandler(&MockFeedService{})

	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	w := httptest.NewRecorder()

	handler.GetFeed(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestFeedHandler_GetFeed_InvalidCursor(t *testing.T) {
	userID := uuid.New()

	handler, _ := NewFeedHandler(&MockFeedService{})

	req := httptest.NewRequest(http.MethodGet, "/feed?cursor=invalid", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.GetFeed(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestFeedHandler_GetFeed_ServiceError(t *testing.T) {
	userID := uuid.New()

	mockService := &MockFeedService{
		GetFeedFn: func(ctx context.Context, uid uuid.UUID, before time.Time, limit int) ([]*post.Post, time.Time, error) {
			return nil, time.Time{}, ErrServiceError
		},
	}

	handler, _ := NewFeedHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.GetFeed(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}


