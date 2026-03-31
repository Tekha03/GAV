package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/profile"
	"social_network/transport/http/middleware"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupProfileHandler(t *testing.T) (*ProfileHandler, *MockProfileService) {
	mockService := new(MockProfileService)

	handler, err := NewProfileHandler(mockService)
	require.NoError(t, err)

	return handler, mockService
}

func withUser(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, middleware.UserIDKey, userID)
}

func withURLParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestProfileHandler_Create_Success(t *testing.T) {
	handler, service := setupProfileHandler(t)

	userID := uuid.New()

	input := profile.CreateProfileInput{
		Name: "John",
	}

	expected := &profile.UserProfile{
		UserID: userID,
		Name:   "John",
	}

	service.
		On("Create", mock.Anything, userID, input).
		Return(expected, nil)

	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(body))
	req = req.WithContext(withUser(req.Context(), userID))

	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	service.AssertExpectations(t)
}

func TestProfileHandler_Create_Unauthorized(t *testing.T) {
	handler, _ := setupProfileHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/profile", nil)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProfileHandler_GetByID_Success(t *testing.T) {
	handler, service := setupProfileHandler(t)

	id := uuid.New()

	expected := &profile.UserProfile{
		UserID: id,
		Name:   "John",
	}

	service.
		On("GetByID", mock.Anything, id).
		Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile/"+id.String(), nil)
	req = withURLParam(req, "id", id.String())

	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProfileHandler_GetByID_InvalidID(t *testing.T) {
	handler, _ := setupProfileHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/profile/invalid", nil)
	req = withURLParam(req, "id", "invalid")

	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProfileHandler_GetByUserID_Success(t *testing.T) {
	handler, service := setupProfileHandler(t)

	userID := uuid.New()

	expected := &profile.UserProfile{
		UserID: userID,
	}

	service.
		On("GetByUserID", mock.Anything, userID).
		Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile/user/"+userID.String(), nil)
	req = withURLParam(req, "userID", userID.String())

	w := httptest.NewRecorder()

	handler.GetByUserID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProfileHandler_Update_Success(t *testing.T) {
	handler, service := setupProfileHandler(t)

	id := uuid.New()

	name := "Updated"
	input := profile.UpdateProfileInput{
		Name: &name,
	}

	service.
		On("Update", mock.Anything, id, input).
		Return(nil)

	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPut, "/profile/"+id.String(), bytes.NewReader(body))
	req = withURLParam(req, "id", id.String())

	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestProfileHandler_Update_InvalidID(t *testing.T) {
	handler, _ := setupProfileHandler(t)

	req := httptest.NewRequest(http.MethodPut, "/profile/invalid", nil)
	req = withURLParam(req, "id", "invalid")

	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProfileHandler_Delete_Success(t *testing.T) {
	handler, service := setupProfileHandler(t)

	id := uuid.New()

	service.
		On("Delete", mock.Anything, id).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/profile/"+id.String(), nil)
	req = withURLParam(req, "id", id.String())

	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestProfileHandler_Delete_InvalidID(t *testing.T) {
	handler, _ := setupProfileHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/profile/invalid", nil)
	req = withURLParam(req, "id", "invalid")

	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
