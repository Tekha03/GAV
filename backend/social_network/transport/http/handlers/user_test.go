package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"social_network/internal/dog"
	"social_network/internal/user"
	"social_network/transport/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserHandler(t *testing.T) *UserHandler {
	mockService := new(MockUserService)
	handler, err := NewUserHandler(mockService)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func TestUserHandler_GetByID_Success(t *testing.T) {
	handler := setupUserHandler(t)
	mockService := handler.service.(*MockUserService)

	userID := uuid.New()
	mockUser := &user.User{Email: "test@example.com"}

	mockService.On("GetByID", mock.Anything, userID).Return(mockUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
	w := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chi.RouteContext(req.Context()).URLParams.Add("id", userID.String())

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), mockUser.Email)
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetByID_InvalidUUID(t *testing.T) {
	handler := setupUserHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/users/bad-id", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chi.RouteContext(req.Context()).URLParams.Add("id", "bad-id")

	w := httptest.NewRecorder()
	handler.GetByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Update_Success(t *testing.T) {
	handler := setupUserHandler(t)
	mockService := handler.service.(*MockUserService)

	userID := uuid.New()
	email := "new@example.com"
	input := user.UpdateUserInput{Email: &email}
	body, _ := json.Marshal(input)

	mockService.On("Update", mock.Anything, userID, input).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chi.RouteContext(req.Context()).URLParams.Add("id", userID.String())

	w := httptest.NewRecorder()
	handler.Update(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_Update_InvalidUUID(t *testing.T) {
	handler := setupUserHandler(t)
	body := []byte(`{"email":"new@example.com"}`)

	req := httptest.NewRequest(http.MethodPut, "/users/bad-id", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chi.RouteContext(req.Context()).URLParams.Add("id", "bad-id")

	w := httptest.NewRecorder()
	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Update_InvalidBody(t *testing.T) {
	handler := setupUserHandler(t)
	userID := uuid.New()

	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader([]byte(`bad json`)))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chi.RouteContext(req.Context()).URLParams.Add("id", userID.String())

	w := httptest.NewRecorder()
	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Delete_Success(t *testing.T) {
	handler := setupUserHandler(t)
	mockService := handler.service.(*MockUserService)

	userID := uuid.New()
	mockService.On("Delete", mock.Anything, userID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+userID.String(), nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chi.RouteContext(req.Context()).URLParams.Add("id", userID.String())

	w := httptest.NewRecorder()
	handler.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_Delete_InvalidUUID(t *testing.T) {
	handler := setupUserHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/users/bad-id", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chi.RouteContext(req.Context()).URLParams.Add("id", "bad-id")

	w := httptest.NewRecorder()
	handler.Delete(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDogHandler_FindDogsNearby_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler, _ := NewUserHandler(mockService)

	userID := uuid.New()
	lat, lon, radius := 55.75, 37.61, 1000.0

	expectedDogs := []*dog.Dog{
		{ID: uuid.New(), Name: "Rex", LocationStatus: 1},
	}

	mockService.On("FindDogsNearby", mock.Anything, userID, lat, lon, radius).
		Return(expectedDogs, nil)

	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=55.75&lon=37.61&radius=1000", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.FindDogsNearby(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDogHandler_FindDogsNearby_InvalidParams(t *testing.T) {
	mockService := new(MockUserService)
	handler, _ := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=bad&lon=37.61&radius=1000", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New())
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.FindDogsNearby(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDogHandler_FindDogsNearby_Unauthorized(t *testing.T) {
	mockService := new(MockUserService)
	handler, _ := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=55.75&lon=37.61&radius=1000", nil)

	rr := httptest.NewRecorder()
	handler.FindDogsNearby(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
