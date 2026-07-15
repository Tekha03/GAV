package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/stats"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type statsHandlerEnv struct {
	handler *StatsHandler
	service *MockStatsService
}

func setupStatsHandler(t *testing.T) *statsHandlerEnv {
	service := &MockStatsService{}
	handler, err := NewStatsHandler(service)
	assert.NoError(t, err)

	return &statsHandlerEnv{
		handler: handler,
		service: service,
	}
}

func TestStatsHandler_UserStats_Success(t *testing.T) {
	env := setupStatsHandler(t)
	userID := uuid.New()
	expected := &stats.UserStats{UserID: userID}

	env.service.On("UserStats", mock.Anything, userID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/stats/users/"+userID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	env.handler.UserStats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got stats.UserStats
	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, expected, &got)
}

func TestStatsHandler_ProfileStats_Success(t *testing.T) {
	env := setupStatsHandler(t)
	userID := uuid.New()
	expected := &stats.ProfileStats{UserID: userID}

	env.service.On("ProfileStats", mock.Anything, userID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/stats/profiles/"+userID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	env.handler.ProfileStats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got stats.ProfileStats
	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, expected, &got)
}

func TestStatsHandler_PostStats_Success(t *testing.T) {
	env := setupStatsHandler(t)
	postID := uuid.New()
	expected := &stats.PostStats{PostID: postID}

	env.service.On("PostStats", mock.Anything, postID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/stats/posts/"+postID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("postID", postID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	env.handler.PostStats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got stats.PostStats
	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, expected, &got)
}

func TestStatsHandler_InvalidUUID(t *testing.T) {
	env := setupStatsHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/stats/users/invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	env.handler.UserStats(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/stats/profiles/invalid", nil)
	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("userID", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	env.handler.ProfileStats(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/stats/posts/invalid", nil)
	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("postID", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	env.handler.PostStats(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
